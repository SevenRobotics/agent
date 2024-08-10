// main package.
package main

import (
	"go_agent/cmd/code-generator/rostoproto"
	"go_agent/utils"
	"sync"
)

var cli struct {
	GoPackage  string `name:"gopackage" help:"Go package name" default:"main"`
	RosPackage string `name:"rospackage" help:"ROS package name" default:"my_package"`
	Path       string `arg:"" help:"path pointing to a ROS message"`
}

var g = rostoproto.NewGen()

func main() {
	gState := utils.GeneratorState{}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go rostoproto.Run(g, &gState, &wg)
	wg.Wait()

	// rmqConfig := config.RMQConfig{
	// 	Username: "seven-admin",
	// 	Password: "123456",
	// 	Host:     "ec2-18-61-25-137.ap-south-2.compute.amazonaws.com:5672",
	// 	Vhost:    "robots",
	// }
	//
	// nodeConfig := goroslib.NodeConf{
	// 	Name:          "telemetry_conductor",
	// 	MasterAddress: "127.0.0.1:11311",
	// }
	//
	// conductor, err := channel.NewConductor(rmqConfig, nodeConfig)
	//
	// if err != nil {
	// 	log.Fatalf("%s", err)
	// }
	//
	// conductor.Start(&gState)

}
