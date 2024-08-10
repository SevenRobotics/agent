package main

import (
	"go_agent/code-generator/rostoproto"
	"go_agent/utils"
	"sync"
)

var g = rostoproto.NewGen()

func main() {
	gState := utils.GeneratorState{}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go rostoproto.Run(g, &gState, &wg)
	wg.Wait()
}
