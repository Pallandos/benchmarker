package main

import (
	"fmt"

	"github.com/pallandos/benchmarker/internal/config"
	"github.com/pallandos/benchmarker/internal/containers"
)

func main() {
	cfg, err := config.LoadConfig("configs/.env")
	if err != nil {
		panic(err)
	}

	stackname := cfg.StackName

	containers, err := containers.ListContainerInfos(stackname)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found containers: %+v\n", containers)
}
