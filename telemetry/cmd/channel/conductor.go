package channel

import (
	"context"
	"errors"
	"fmt"
	"go_agent/config"
	"go_agent/iface"
	"go_agent/publishers/rmq"
	"go_agent/utils"
	"log"
	"time"

	// "log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/bluenviron/goroslib/v2/pkg/apimaster"
)

type MessageInfo struct {
	Name    string
	Package string
}

type TopicInfo map[string]MessageInfo

type State struct {
	Topics TopicInfo //map of topic name to message type

	Builders map[string]iface.Builder

	Pipelines map[string]iface.Pipeline

	Configs map[string]*config.RRPipelineConfig

	ValidTopics map[string]struct{}

	UserRequestedTopics map[string]struct{}
}

type Conductor interface {
	RunPipelines(*sync.WaitGroup)
	BuildPipelines() error
	Start(genState *utils.GeneratorState, userTopics []string) error
}

type conductor struct {

	//active internal state of the topics and their pipelines/builders maintained by the conductor
	internalState State

	//error channels per pipeline
	errorChannels map[string]chan error

	//wait group for pipelines
	waitGroup *sync.WaitGroup

	//node used by the conductor to query info regarding the topics and nodes
	client     *apimaster.Client
	nodeConfig *config.RosNodeConfig

	rmqConfig config.RMQConfig
	genState  *utils.GeneratorState

	ticker      *time.Ticker
	builderutil utils.BuilderFinder
}

const dependencyRetryInterval = 5 * time.Second

const dependencyHealthClientName = "__telemetry_dependency_health"

func NewConductor(rmqConf config.RMQConfig, nodeConfig config.RosNodeConfig) (Conductor, error) {

	builder_util := utils.NewBuilder("ros-rmq", nil)
	if !builder_util.Generated {
		return nil, fmt.Errorf("builder util not generated or assigned yet. Do not call build code during generation")
	}

	return &conductor{
		rmqConfig:   rmqConf,
		nodeConfig:  &nodeConfig,
		client:      apimaster.NewClient(nodeConfig.Address, nodeConfig.Name, &http.Client{}),
		builderutil: builder_util,
	}, nil
}

func normalizeAgentID(agentID string) string {
	normalized := strings.ToLower(strings.TrimSpace(agentID))
	if normalized == "" {
		return ""
	}

	normalized = strings.ReplaceAll(normalized, "_", "")
	normalized = strings.ReplaceAll(normalized, "-", "")

	if strings.HasPrefix(normalized, "amr") {
		suffix := strings.TrimPrefix(normalized, "amr")
		suffix = strings.TrimPrefix(suffix, ".")
		if id, err := strconv.Atoi(suffix); err == nil {
			return fmt.Sprintf("amr.%03d", id)
		}
	}

	return normalized
}

func exchangeName(agentID string) string {
	return normalizeAgentID(agentID) + ".exchange"
}

func topicRouteName(topic string) string {
	return strings.Trim(strings.TrimSpace(topic), ".")
}

func routingKey(agentID, topic string) string {
	return normalizeAgentID(agentID) + "." + topicRouteName(topic)
}

func queueName(agentID, topic string) string {
	return routingKey(agentID, topic) + ".q"
}

func (c *conductor) Start(genState *utils.GeneratorState, userTopics []string) error {

	c.genState = genState

	c.errorChannels = map[string]chan error{}
	c.errorChannels["self"] = make(chan error, 1)

	c.internalState.Builders = map[string]iface.Builder{}
	c.internalState.Pipelines = map[string]iface.Pipeline{}
	c.internalState.Configs = map[string]*config.RRPipelineConfig{}
	c.internalState.UserRequestedTopics = map[string]struct{}{}

	for _, topic := range userTopics {
		c.internalState.UserRequestedTopics[topic] = struct{}{}
	}

	c.waitGroup = &sync.WaitGroup{}

	for {
		c.waitForDependencies()
		c.resetTopicState()

		if err := c.rebuildConfiguredPipelines(); err != nil {
			log.Printf("Error rebuilding configured pipelines: %v", err)
			c.shutdownAllPipelines()
			time.Sleep(dependencyRetryInterval)
			continue
		}

		c.drainDependencyErrors()
		c.RunPipelines(c.waitGroup)

		err := c.monitorDependencies()
		log.Printf("Dependency failure detected, shutting down pipelines: %v", err)
		c.shutdownAllPipelines()
	}
}

func (c *conductor) isRosMasterAvailable() bool {
	_, err := c.client.GetPublishedTopics("")
	return err == nil
}

func (c *conductor) isRabbitMQAvailable() bool {
	conn, err := rmq.NewRabbitMQ(c.rmqConfig)
	if err != nil {
		return false
	}
	if !conn.IsConnected() {
		return false
	}

	_, err = conn.NewClient(dependencyHealthClientName)
	return err == nil
}

func (c *conductor) waitForDependencies() {
	for {
		rosAvailable := c.isRosMasterAvailable()
		rmqAvailable := c.isRabbitMQAvailable()
		if rosAvailable && rmqAvailable {
			log.Printf("ROS Master and RabbitMQ are online")
			return
		}

		if !rosAvailable {
			log.Printf("ROS Master not available at %s, retrying in %s...", c.nodeConfig.Address, dependencyRetryInterval)
		}
		if !rmqAvailable {
			log.Printf("RabbitMQ not available, retrying in %s...", dependencyRetryInterval)
		}
		time.Sleep(dependencyRetryInterval)
	}
}

func (c *conductor) monitorDependencies() error {
	ticker := time.NewTicker(dependencyRetryInterval)
	defer ticker.Stop()

	for {
		select {
		case err := <-c.errorChannels["self"]:
			return err
		case <-ticker.C:
			if !c.isRosMasterAvailable() {
				return fmt.Errorf("ROS Master not available at %s", c.nodeConfig.Address)
			}
			if !c.isRabbitMQAvailable() {
				return fmt.Errorf("%w: RabbitMQ health check failed", rmq.ErrRabbitMQUnavailable)
			}
		}
	}
}

func (c *conductor) drainDependencyErrors() {
	for {
		select {
		case <-c.errorChannels["self"]:
		default:
			return
		}
	}
}

func (c *conductor) rebuildConfiguredPipelines() error {
	if !c.isRosMasterAvailable() {
		return fmt.Errorf("ROS Master not available at %s", c.nodeConfig.Address)
	}
	if !c.isRabbitMQAvailable() {
		return fmt.Errorf("%w: RabbitMQ not available", rmq.ErrRabbitMQUnavailable)
	}

	if err := c.LoadTopicInfo(); err != nil {
		return err
	}
	if err := c.ConfigureBuilders(); err != nil {
		return err
	}
	if err := c.BuildPipelines(); err != nil {
		return err
	}
	return nil
}

func (c *conductor) shutdownAllPipelines() {
	for name, pipeline := range c.internalState.Pipelines {
		if pipeline.IsActive() {
			log.Printf("Shutting down pipeline %s\n", name)
			pipeline.Shutdown()
		}
	}
}

func (c *conductor) resetTopicState() {
	c.internalState.Topics = TopicInfo{}
	c.internalState.ValidTopics = map[string]struct{}{}
	c.internalState.Builders = map[string]iface.Builder{}
	c.internalState.Pipelines = map[string]iface.Pipeline{}
	c.internalState.Configs = map[string]*config.RRPipelineConfig{}
}

func (c *conductor) LoadTopicInfo() error {
	if c.client == nil {
		return fmt.Errorf("xmlrpc client not set")
	}

	topics, err := c.client.GetPublishedTopics("")
	if err != nil {
		return fmt.Errorf("Could not get topics from Ros Master %s: %v", c.nodeConfig.Address, err)
	}

	return c.loadTopicInfo(topics)

}

func (c *conductor) loadTopicInfo(topics [][]string) error {

	if c.internalState.Topics == nil {
		c.internalState.Topics = TopicInfo{}
	}

	if c.internalState.ValidTopics == nil {
		c.internalState.ValidTopics = map[string]struct{}{}
	}

	for _, message := range topics {
		// Only process topics that were explicitly requested by the user
		if _, ok := c.internalState.UserRequestedTopics[message[0]]; !ok {
			continue
		}

		tmp := strings.Split(message[1], "/")
		info := MessageInfo{
			Name:    tmp[1],
			Package: tmp[0],
		}
		t := strings.Split(message[0], "/")
		if len(t) > 2 {
			topicName := strings.ReplaceAll(message[0], "/", ".")
			c.internalState.Topics[topicName[1:]] = info
		} else {
			c.internalState.Topics[t[1]] = info
		}
	}

	for k, v := range c.internalState.Topics {
		if _, ok := c.internalState.UserRequestedTopics["/"+strings.ReplaceAll(k, ".", "/")]; ok {
			if t, ok := c.genState.RosMsgPkgs[v.Package]; ok {
				if _, ok := t[v.Name]; ok {
					c.internalState.ValidTopics[k] = struct{}{}
					continue
				}
			}
		}
	}

	return nil

}

func (c *conductor) ConfigureBuilders() error {
	agentID := normalizeAgentID(c.nodeConfig.AgentID)
	if agentID == "" {
		return fmt.Errorf("agent_id is required in telemetry node config")
	}

	for k := range c.internalState.ValidTopics {
		if info, ok := c.internalState.Topics[k]; ok {
			var err error
			if _, ok := c.internalState.Builders[k]; !ok {
				c.internalState.Builders[k], err = c.builderutil.GetBuilderFromName(info.Name)
				if err != nil {
					return fmt.Errorf("Error retrieving builder for %s:%v", info.Name, err)
				}
				log.Printf("Added builder for %s %s: %v", k, "/"+strings.ReplaceAll(k, ".", "/"), c.internalState.Builders[k])

				c.internalState.Configs[k] = &config.RRPipelineConfig{
					RMQConnConfig: c.rmqConfig,
					RMQPubConfig: config.RMQClientConfig{
						Exchange:   exchangeName(agentID),
						Topic:      queueName(agentID, k),
						RoutingKey: routingKey(agentID, k),
						Durable:    true,
						Autodelete: false,
						Ctx:        context.Background(),
					},
					SubConfig: config.RosSubscriberConfig{
						Node: config.RosNodeConfig{
							Name:    strings.ReplaceAll(k, ".", "_") + "_node",
							Address: c.nodeConfig.Address,
						},
						Topic: "/" + strings.ReplaceAll(k, ".", "/"),
						Name:  strings.ReplaceAll(k, ".", "_") + "_sub",
					},
					Name: k,
				}
			}
		}
	}
	return nil
}

func (c *conductor) BuildPipelines() error {
	var buildErr error
	for name, builder := range c.internalState.Builders {
		conf, ok := c.internalState.Configs[name]
		if !ok {
			log.Printf("Pipeline Config not found for %s, skipping\n", name)
			continue
		}

		if pipe, ok := c.internalState.Pipelines[name]; ok {
			if pipe.IsActive() {
				log.Printf("Pipeline %s is active, will not rebuild\n", pipe.Name())
				continue
			} else {
				log.Printf("Pipeline %s is inactive\n", pipe.Name())
			}
		}

		p, err := builder.BuildPipeline(*conf)
		if err != nil {
			log.Printf("Failed to build pipeline for %s: %v\n", name, err)
			buildErr = errors.Join(buildErr, fmt.Errorf("build pipeline %s: %w", name, err))
			continue
		}
		c.internalState.Pipelines[name] = p
		log.Printf("Pipeline created for %s\n", name)
		c.errorChannels[name] = c.internalState.Pipelines[name].GetErrorStream()
	}
	return buildErr
}

func (c *conductor) hasPendingPipelines() bool {
	for name := range c.internalState.Builders {
		if _, ok := c.internalState.Configs[name]; !ok {
			continue
		}

		pipeline, ok := c.internalState.Pipelines[name]
		if !ok || !pipeline.IsActive() {
			return true
		}
	}
	return false
}

func (c *conductor) RunPipelines(wg *sync.WaitGroup) {
	for _, pipeline := range c.internalState.Pipelines {
		if !pipeline.IsActive() {
			errCh := pipeline.GetErrorStream()
			c.errorChannels[pipeline.Name()] = errCh
			c.waitGroup.Add(1)
			go func(name string, p iface.Pipeline, ch chan error) {
				defer c.waitGroup.Done()
				for {
					select {
					case err, ok := <-ch:
						if !ok {
							return
						}
						log.Printf("Error on %s: %v\n", name, err)
						if errors.Is(err, rmq.ErrRabbitMQUnavailable) || strings.Contains(err.Error(), "failed to initialise subscriber") {
							select {
							case c.errorChannels["self"] <- err:
							default:
							}
						}
					case <-time.After(1 * time.Second):
						if !p.IsActive() {
							return
						}
					}
				}
			}(pipeline.Name(), pipeline, errCh)
			c.waitGroup.Add(1)
			go pipeline.Start(c.waitGroup)
		}
	}
}
