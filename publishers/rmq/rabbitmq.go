package rmq

import (
	"context"
	"fmt"
	"go_agent/config"
	"log"
	"net"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var lock = &sync.Mutex{}

type RabbitMQMaster interface {
	Connect() RabbitMQMaster
	HasClient(name string) bool
	NewClient(name string) (RabbitClient, error)
	RemoveClient(name string) error
	Close() error
}

type RabbitClient interface {
	Close() error
	NewExchangeDeclare(exchangeName, kind string, durable, autodelete bool) error
	NewQueueDeclare(queueName string, durable, autodelete bool) error
	CreateBinding(name string, binding string, exchange string) error
	Send(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error
	Receive(ctx context.Context, queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error)
}

type rabbitClient struct {
	// rule of thumb is to use a single connection per app and spawn channels for every task
	mu   sync.RWMutex
	conn *amqp.Connection // a tcp connection used by the client
	ch   *amqp.Channel    // a multiplexed connection over the tcp connection i.e, conn
}

var rabbitMQSingleton *rabbitMQ

type rabbitMQ struct {
	mu      sync.RWMutex
	conn    *amqp.Connection
	conf    config.RMQConfig
	clients map[string]*rabbitClient
}

func (r *rabbitMQ) Connect() (*rabbitMQ, error) {
	err := r.dial()
	if err != nil {
		return nil, err
	}
	go r.monitorConnection()
	return r, nil
}

func (r *rabbitMQ) dial() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	conn, err := amqp.DialConfig(fmt.Sprintf("amqp://%s:%s@%s/%s",
		r.conf.Username, r.conf.Password, r.conf.Host, r.conf.Vhost),
		amqp.Config{
			Heartbeat: 10 * time.Second,
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, 5*time.Second)
			},
		},
	)

	if err != nil {
		return err
	}

	r.conn = conn
	return nil
}

func (r *rabbitMQ) monitorConnection() {
	for {
		closeCh := make(chan *amqp.Error, 1)
		r.mu.RLock()
		if r.conn == nil {
			r.mu.RUnlock()
			return
		}
		r.conn.NotifyClose(closeCh)
		r.mu.RUnlock()

		err := <-closeCh
		if err == nil {
			// Intentional close
			return
		}

		log.Printf("RMQ connection lost: %v. Attempting to reconnect...", err)
		r.reconnectWithBackoff()
		r.recreateClientChannels()
	}
}

func (r *rabbitMQ) reconnectWithBackoff() {
	backoff := 1 * time.Second
	maxBackoff := 30 * time.Second

	for {
		log.Printf("RMQ reconnecting in %v...", backoff)
		time.Sleep(backoff)

		err := r.dial()
		if err == nil {
			log.Printf("RMQ successfully reconnected")
			return
		}

		log.Printf("RMQ reconnection failed: %v", err)
		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}
	}
}

func (r *rabbitMQ) recreateClientChannels() {
	r.mu.Lock()
	clients := make(map[string]*rabbitClient)
	for k, v := range r.clients {
		clients[k] = v
	}
	r.mu.Unlock()

	for name, client := range clients {
		r.recreateClientChannel(name, client)
	}
}

func (r *rabbitMQ) recreateClientChannel(name string, client *rabbitClient) {
	r.mu.RLock()
	conn := r.conn
	r.mu.RUnlock()

	if conn == nil || conn.IsClosed() {
		return
	}

	client.mu.RLock()
	// If the channel is already open and using the current connection, skip recreation
	if client.ch != nil && !client.ch.IsClosed() && client.conn == conn {
		client.mu.RUnlock()
		return
	}
	client.mu.RUnlock()

	log.Printf("Recreating channel for RMQ client: %s", name)
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to recreate channel for %s: %v. Will retry on next send.", name, err)
		return
	}

	client.mu.Lock()
	client.conn = conn
	client.ch = ch
	client.mu.Unlock()

	// Start monitoring the new channel
	go r.monitorChannel(name, client)
}

func (r *rabbitMQ) monitorChannel(name string, client *rabbitClient) {
	closeCh := make(chan *amqp.Error, 1)
	client.mu.RLock()
	if client.ch == nil {
		client.mu.RUnlock()
		return
	}
	client.ch.NotifyClose(closeCh)
	client.mu.RUnlock()

	err := <-closeCh
	if err == nil {
		return // Intentional close
	}

	log.Printf("RMQ channel lost for client %s: %v. Attempting to recreate channel...", name, err)
	
	// Check if the connection is still alive. If it is, just recreate the channel.
	// If the connection is dead, the connection monitor will handle it.
	r.mu.RLock()
	connAlive := r.conn != nil && !r.conn.IsClosed()
	r.mu.RUnlock()

	if connAlive {
		r.recreateClientChannel(name, client)
	}
}

func createRabbitMQ(conf config.RMQConfig) (*rabbitMQ, error) {

	if conf.Username == "" {
		return nil, fmt.Errorf("Empty username provided")
	}

	if conf.Password == "" {
		return nil, fmt.Errorf("Empty password provided")
	}

	if conf.Host == "" {
		return nil, fmt.Errorf("host cannot be an empty string")
	}

	if conf.Vhost == "" {
		return nil, fmt.Errorf("vhost cannot be an empty string")
	}

	r := &rabbitMQ{
		conf:    conf,
		clients: map[string]*rabbitClient{},
	}

	_, err := r.Connect()

	if err != nil {
		return nil, fmt.Errorf("Failed to connect with %v: %v", r.conf, err)
	}

	return r, nil

}

func NewRabbitMQ(conf config.RMQConfig) (*rabbitMQ, error) {

	if rabbitMQSingleton == nil {
		lock.Lock()
		defer lock.Unlock()
		r, err := createRabbitMQ(conf)
		if err != nil {
			return nil, fmt.Errorf("RMQ Creation error: %v", err)
		}
		rabbitMQSingleton = r
	}
	return rabbitMQSingleton, nil
}

func (r *rabbitMQ) HasClient(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if _, ok := r.clients[name]; !ok {
		return false
	}
	return true
}

func (r *rabbitMQ) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

func (r *rabbitMQ) NewClient(name string) (*rabbitClient, error) {
	r.mu.Lock()
	client, exists := r.clients[name]
	if exists {
		r.mu.Unlock()
		// Ensure the existing client has a valid channel
		r.recreateClientChannel(name, client)
		return client, nil
	}

	if r.conn == nil {
		r.mu.Unlock()
		return nil, fmt.Errorf("Cannot add clients to an uninitialized RMQ Connection")
	}

	ch, err := r.conn.Channel()
	if err != nil {
		r.mu.Unlock()
		return nil, err
	}

	client = &rabbitClient{
		conn: r.conn,
		ch:   ch,
	}
	r.clients[name] = client
	r.mu.Unlock()

	// Start monitoring the channel
	go r.monitorChannel(name, client)

	return client, nil
}

func (r *rabbitMQ) RemoveClient(name string) error {
	r.mu.Lock()
	client, ok := r.clients[name]
	if !ok {
		r.mu.Unlock()
		return fmt.Errorf("Client %s does not exist", name)
	}
	delete(r.clients, name)
	r.mu.Unlock()

	return client.Close()
}

func (rc *rabbitClient) Close() error {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	if rc.ch != nil {
		return rc.ch.Close()
	}
	return nil
}

func (rc *rabbitClient) NewExchangeDeclare(exchangeName, kind string, durable, autodelete bool) error {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	err := rc.ch.ExchangeDeclare(exchangeName, kind, durable, autodelete, false, false, nil)
	return err
}

func (rc *rabbitClient) NewQueueDeclare(queueName string, durable, autodelete bool) error {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	_, err := rc.ch.QueueDeclare(queueName, durable, autodelete, false, false, nil)
	return err
}

func (rc *rabbitClient) CreateBinding(name string, binding string, exchange string) error {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.ch.QueueBind(name, binding, exchange, false, nil)
}

func (rc *rabbitClient) Send(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.ch.PublishWithContext(ctx, exchange, routingKey, true, false, options)
}

func (rc *rabbitClient) Receive(ctx context.Context, queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.ch.ConsumeWithContext(ctx, queue, consumer, autoAck, false, false, false, nil)
}
