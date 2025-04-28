package main

import (
	"go_agent/config"
	"go_agent/publishers/rmq"
	"go_agent/telemetry/cmd/channel"
	"go_agent/telemetry/gengo/ros/converter"
	"go_agent/utils"
	"io/fs"
	"log"
	"context"
	"os"
	"path/filepath"
	"runtime"
	"fmt"
	"strings"
	"sync"
	// "github.com/rabbitmq/amqp091-go"

	"gopkg.in/yaml.v2"
)

var genState utils.GeneratorState

type GeneratedFilesInfo struct {
	ProtoDir     string
	GoDir        string
	ConverterDir string
	ConfigDir    string
}

func NewGen() *GeneratedFilesInfo {
	defaultOutputDir := "genproto/ros/"
	goOutputDir := "gengo/ros/"
	converterDir := "gengo/ros/converter/"
	configDir := "../config/"
	return &GeneratedFilesInfo{
		GoDir:        goOutputDir,
		ConverterDir: converterDir,
		ProtoDir:     defaultOutputDir,
		ConfigDir:    configDir,
	}
}

func main() {
	g := NewGen()

	//scan and fill generated state
	genState = utils.GeneratorState{}
	genState.RosMsgPkgs = map[string]map[string]interface{}{}

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	var wg sync.WaitGroup

	rosGoDir := filepath.Join(basepath, g.ProtoDir)
	var pkgname string
	filepath.WalkDir(rosGoDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			pkgname = d.Name()
		}
		if t, ok := genState.RosMsgPkgs[pkgname]; !ok {
			genState.RosMsgPkgs[pkgname] = map[string]interface{}{}
		} else {
			msgname := d.Name()
			t[strings.TrimSuffix(msgname, ".pb.go")] = struct{}{}
		}
		return nil
	})

	if len(genState.RosMsgPkgs) == 0 {
		log.Fatalf("Failed to collect generated output info")
	}

	if !converter.AssignBuilder() {
		log.Fatalf("Pipeline builder utility not assigned")
	}

	config_path := filepath.Join(basepath, g.ConfigDir, "rmq_config.yml")
	f, err := os.Open(config_path)
	if err != nil {
		log.Fatalf("RMQ Configuration not found @ %s: %v", config_path, err)
	}
	defer f.Close()

	var rmq_config config.RMQConfig
	decoder := yaml.NewDecoder(f)

	err = decoder.Decode(&rmq_config)
	if err != nil {
		log.Fatalf("Error decoding RMQ Config from %s: %v", config_path, err)
	}

	config_path = filepath.Join(basepath, g.ConfigDir, "telemetry_node.yml")
	nf, err := os.Open(config_path)
	if err != nil {
		log.Fatalf("Node Configuration not found @ %s: %v", config_path, err)
	}
	defer nf.Close()

	var node_config config.RosNodeConfig
	nodeDecoder := yaml.NewDecoder(nf)



	err = nodeDecoder.Decode(&node_config)
	if err != nil {
		log.Fatalf("Error decoding Node Config from %s: %v", config_path, err)
	}


	topicList := []string{"/odom_with_amcl","/scan"}

	conn, err := rmq.NewRabbitMQ(rmq_config)

	if err != nil {
		fmt.Errorf("Failed to connect to RMQ Server:%v", err)
	}

	client, err := conn.NewClient("tasksubscriber")

	if err != nil {
	    fmt.Errorf("Failed to create RMQ Client : %v", err)
	}

	// Create subscriber configuration
	subscriberConfig := config.RMQClientConfig{
		Exchange:   "robot_exchange",
		Topic:      "task_queue",
		RoutingKey: "task_key",
		Durable:    true,
		Autodelete: false,
		Ctx:        context.Background(),
	}

	subscriber, err := rmq.NewRMQSubscriber[any](subscriberConfig, client)
	if err != nil {
		log.Fatalf("Error decoding Node Config from: %v", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		
		log.Printf("Starting to listen for messages on queue 'task' with routing key 'task_key'")
		deliveries, err := subscriber.Receive()
			if err != nil {
				log.Printf("Error receiving messages: %v", err)
				return
			}
		
		for {
			
			for msg := range deliveries {
				log.Printf("Received message from %s: %s", msg.RoutingKey, string(msg.Body))
				
				// Acknowledge the message to remove it from the queue
				if err := msg.Ack(false); err != nil {
					log.Printf("Error acknowledging message: %v", err)
				}
			}
		}
		
		log.Printf("Subscriber stopped")
	}()

	conductor, err := channel.NewConductor(rmq_config, node_config)
	if err != nil {
		log.Fatalf("New Conductor Could not be created: %v", err)
	}

	err = conductor.Start(&genState, topicList)
	if err != nil {
		log.Fatalf("Conductor Failed: %v", err)
	}


	
	wg.Wait()
}
