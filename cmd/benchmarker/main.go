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
	err = logger.InitLogger("benchmarker.log", cfg.LogPath)
	if err != nil {
		panic(err)
	}

	containers, err := containers.ListContainerInfos(stackname)
	if err != nil {
		logger.Log.Error("Failed to list containers: ", err)
		panic(err)
	}

	logger.Log.Info("Found containers: ", len(containers))
}
