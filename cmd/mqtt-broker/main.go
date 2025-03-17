package main

import (
	"github.com/ApplyLogic/mqtt-broker/app"
	"github.com/ApplyLogic/mqtt-broker/config"
	"log"
	"os"
	"runtime/pprof"
)

func main() {

	cpuPprof, _ := os.Create("cpu.pprof")
	defer cpuPprof.Close()
	_ = pprof.StartCPUProfile(cpuPprof)
	defer pprof.StopCPUProfile()

	blockPprof, _ := os.Create("block.pprof")
	defer blockPprof.Close()
	_ = pprof.Lookup("block").WriteTo(blockPprof, 0)

	memoryPprof, _ := os.Create("memory.pprof")
	defer memoryPprof.Close()
	_ = pprof.WriteHeapProfile(memoryPprof)

	cfg, err := config.LoanConfig()
	if err != nil {
		log.Fatal(err)
	}
	application := &app.App{}
	application.Initialize(cfg)
	application.Start()
}
