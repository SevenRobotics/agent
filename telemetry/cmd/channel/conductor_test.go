package channel

import (
	"go_agent/config"
	"go_agent/iface"
	"sync"
	"testing"
)

type testPipeline struct {
	active bool
	name   string
	errCh  chan error
}

func (p *testPipeline) Shutdown()                  {}
func (p *testPipeline) Start(*sync.WaitGroup)      {}
func (p *testPipeline) GetErrorStream() chan error { return p.errCh }
func (p *testPipeline) Name() string               { return p.name }
func (p *testPipeline) IsActive() bool             { return p.active }
func (p *testPipeline) Deactivate()                { p.active = false }

func TestConductorHasPendingPipelines(t *testing.T) {
	tests := []struct {
		name      string
		state     State
		wantRetry bool
	}{
		{
			name: "missing configured pipeline",
			state: State{
				Builders: map[string]iface.Builder{
					"odom": nil,
				},
				Configs: map[string]*config.RRPipelineConfig{
					"odom": {},
				},
				Pipelines: map[string]iface.Pipeline{},
			},
			wantRetry: true,
		},
		{
			name: "inactive configured pipeline",
			state: State{
				Builders: map[string]iface.Builder{
					"odom": nil,
				},
				Configs: map[string]*config.RRPipelineConfig{
					"odom": {},
				},
				Pipelines: map[string]iface.Pipeline{
					"odom": &testPipeline{name: "odom", active: false, errCh: make(chan error)},
				},
			},
			wantRetry: true,
		},
		{
			name: "active configured pipeline",
			state: State{
				Builders: map[string]iface.Builder{
					"odom": nil,
				},
				Configs: map[string]*config.RRPipelineConfig{
					"odom": {},
				},
				Pipelines: map[string]iface.Pipeline{
					"odom": &testPipeline{name: "odom", active: true, errCh: make(chan error)},
				},
			},
			wantRetry: false,
		},
		{
			name: "builder without config is ignored",
			state: State{
				Builders: map[string]iface.Builder{
					"odom": nil,
				},
				Configs:   map[string]*config.RRPipelineConfig{},
				Pipelines: map[string]iface.Pipeline{},
			},
			wantRetry: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &conductor{internalState: tt.state}
			if got := c.hasPendingPipelines(); got != tt.wantRetry {
				t.Fatalf("hasPendingPipelines() = %v, want %v", got, tt.wantRetry)
			}
		})
	}
}
