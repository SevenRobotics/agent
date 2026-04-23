package config

import (
	"context"
)

type RosNodeConfig struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
	AgentID string `yaml:"agent_id"`
}

type RosSubscriberConfig struct {
	Node  RosNodeConfig
	Topic string
	Name  string //subscriber name; will be same as channel name
}

type RMQConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Vhost    string `yaml:"vhost"`
}

type RMQClientConfig struct {
	Exchange   string
	Topic      string
	RoutingKey string
	Durable    bool
	Autodelete bool
	Ctx        context.Context
}

// Ros Rmq Pipeline Config
type RRPipelineConfig struct {
	SubConfig     RosSubscriberConfig
	RMQConnConfig RMQConfig
	RMQPubConfig  RMQClientConfig
	Name          string //name of the channel
}
