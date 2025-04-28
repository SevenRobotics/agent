package main

import (
	"go_agent/config"
	"go_agent/telemetry/cmd/channel"
	"go_agent/telemetry/gengo/ros/converter"
	"go_agent/utils"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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

	topicList := []string{"/cmd_vel","/scan"}

	conductor, err := channel.NewConductor(rmq_config, node_config)
	if err != nil {
		log.Fatalf("New Conductor Could not be created: %v", err)
	}

	err = conductor.Start(&genState, topicList)
	if err != nil {
		log.Fatalf("Conductor Failed: %v", err)
	}
}
