package iface

import "go_agent/config"

type Builder interface {
	BuildPipeline(config.RRPipelineConfig) (Pipeline, error)
}
