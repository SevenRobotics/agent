package channel

import (
	"fmt"
	"go_agent/config"
	"go_agent/iface"
)

type BuilderUtil[S any, P any] struct {
	MsgConverter func(in S) (P, error)
	Serializer   func(msg P) ([]byte, error)
}

func (b *BuilderUtil[S, P]) BuildPipeline(conf config.RRPipelineConfig) (iface.Pipeline, error) {
	fmt.Printf("Builder called for \n")
	p, err := NewPipeline(conf, b.MsgConverter, b.Serializer)

	if err != nil {
		return nil, fmt.Errorf("Failure in creating pipeline %s: %v", conf.Name, err)
	}

	return p, nil
}
