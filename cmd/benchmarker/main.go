package main

import (
	"github.com/pallandos/benchmarker/internal/containers"
	"github.com/pallandos/benchmarker/internal/utils/config"
	"github.com/pallandos/benchmarker/internal/utils/logger"
)

func main() {
	cfg, err := config.LoadConfig("configs/.env")
	if err != nil {
		panic(err)
	}

	stackname := cfg.StackName
	mainlogger, err := logger.InitLogger("benchmarker.log", cfg.LogPath)
	if err != nil {
		panic(err)
	}

	containers, err := containers.ListContainerInfos(stackname)
	if err != nil {
		mainlogger.Error("Failed to list containers: ", err)
		panic(err)
	}

	mainlogger.Info("Found containers: ", len(containers))
}
