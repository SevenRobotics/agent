package channel

import (
	"fmt"
	"go_agent/config"
	"go_agent/iface"
	"go_agent/utils"
	"log"
	"time"

	// "log"
	"net/http"
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

	ValidTopics map[string]bool
}

type Conductor interface {
	RunPipelines() error
	BuildPipelines() error
	CheckForNewTopics(chan<- [][]string, chan int)
	Start(genState *utils.GeneratorState) error
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
	//used by conductor discovery goroutine to publish new TopicInfo
	topicDiscoveryChannel chan TopicInfo

	genState *utils.GeneratorState

	ticker *time.Ticker
}

func NewConductor(rmqConf config.RMQConfig, nodeConfig config.RosNodeConfig) (Conductor, error) {
	return &conductor{
		rmqConfig:  rmqConf,
		nodeConfig: &nodeConfig,
		client:     apimaster.NewClient(nodeConfig.Address, nodeConfig.Name, &http.Client{}),
	}, nil
}

func (c *conductor) RunPipelines() error {
	return nil
}

func (c *conductor) BuildPipelines() error {
	return nil
}

func (c *conductor) CheckForNewTopics(topicStream chan<- [][]string, done chan int) {
	c.ticker = time.NewTicker(100 * time.Millisecond)
	defer c.ticker.Stop()
	defer c.waitGroup.Done()

	for {
		select {
		case <-c.ticker.C:
			topics, err := c.client.GetPublishedTopics("")
			if err != nil {
				c.errorChannels["self"] <- err
			}
			newtopics := [][]string{}
			for _, topic := range topics {
				t := strings.Split(topic[0], "/")
				var topicName string
				if len(t) > 2 {
					topicName = strings.ReplaceAll(topic[0], "/", ".")
					topicName = topicName[1:]
				} else {
					topicName = t[1]
				}

				if _, ok := c.internalState.Topics[topicName]; !ok {
					newtopics = append(newtopics, topic)
				}
			}

			if len(newtopics) > 0 {
				topicStream <- newtopics
			}
		case <-done:
			return
		}
	}
}

func (c *conductor) Start(genState *utils.GeneratorState) error {

	c.genState = genState

	c.errorChannels = map[string]chan error{}
	c.errorChannels["self"] = make(chan error)

	c.internalState.Builders = map[string]iface.Builder{}
	c.internalState.Pipelines = map[string]iface.Pipeline{}
	c.internalState.Configs = map[string]*config.RRPipelineConfig{}

	err := c.LoadTopicInfo()

	if err != nil {
		return err
	}

	builder_util := utils.NewBuilder("ros-rmq", nil)
	if !builder_util.Generated {
		return fmt.Errorf("builder util not generated or assigned yet. Do not call build code during generation")
	}

	for k, v := range c.internalState.ValidTopics {
		if v {
			if info, ok := c.internalState.Topics[k]; ok {
				var err error
				c.internalState.Builders[k], err = builder_util.GetBuilderFromName(info.Name)
				if err != nil {
					return fmt.Errorf("Error retrieving builder for %s:%v", info.Name, err)
				}
				log.Printf("Added builder for %s %s: %v", k, "/"+strings.ReplaceAll(k, ".", "/"), c.internalState.Builders[k])
				c.internalState.Configs[k] = &config.RRPipelineConfig{
					RMQConnConfig: c.rmqConfig,
					RMQPubConfig: config.RMQClientConfig{
						Exchange:   "robot1",
						Topic:      k,
						RoutingKey: strings.Join([]string{info.Package, info.Name}, ".") + ".*",
						Durable:    true,
						Autodelete: false,
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

	for name, builder := range c.internalState.Builders {
		if conf, ok := c.internalState.Configs[name]; ok {
			c.internalState.Pipelines[name], err = builder.BuildPipeline(*conf)
			if err != nil {
				return fmt.Errorf("Failed to build pipeline for %s: %v", name, err)
			}
		} else {
			return fmt.Errorf("Pipeline Config not found for %s", name)
		}
		log.Printf("Pipeline created for %s\n", name)
		c.errorChannels[name] = c.internalState.Pipelines[name].GetErrorStream()
	}

	//run error channel listener goroutines here
	for name, errCh := range c.errorChannels {
		go func(errCh chan error) {
			for {
				select {
				case err := <-errCh:
					log.Printf("Error on %s: %v\n", name, err)
				}
			}
		}(errCh)
	}

	topicStream := make(chan [][]string)
	done := make(chan int)
	c.waitGroup = &sync.WaitGroup{}
	c.waitGroup.Add(1)
	go func(topicS chan [][]string, done chan int, conductor *conductor) {
		defer c.waitGroup.Done()
		for {
			select {
			case topics := <-topicS:
				log.Printf("New topic discovered %v\n", topics)
				c.LoadTopicInfoFrom(topics)
			case err := <-c.errorChannels["self"]:
				log.Printf("Error scanning for new topics: %v\n", err)
			case <-done:
				return
			}
		}
	}(topicStream, done, c)
	c.waitGroup.Add(1)

	go c.CheckForNewTopics(topicStream, done)

	//start pipelines here
	for _, pipeline := range c.internalState.Pipelines {
		c.waitGroup.Add(1)
		go pipeline.Start(c.waitGroup)
	}

	c.waitGroup.Wait()
	return nil
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

func (c *conductor) LoadTopicInfoFrom(topics [][]string) error {
	return c.loadTopicInfo(topics)
}

func (c *conductor) loadTopicInfo(topics [][]string) error {

	if c.internalState.Topics == nil {
		c.internalState.Topics = TopicInfo{}
	}

	if c.internalState.ValidTopics == nil {
		c.internalState.ValidTopics = map[string]bool{}
	}

	for _, message := range topics {
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
		if t, ok := c.genState.RosMsgPkgs[v.Package]; ok {
			if _, ok := t[v.Name]; !ok {
				c.internalState.ValidTopics[k] = false
				continue
			}
			c.internalState.ValidTopics[k] = true
		} else {
			c.internalState.ValidTopics[k] = false
		}
	}

	return nil

}
