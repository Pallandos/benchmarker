package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pallandos/benchmarker/internal/containers"
	"github.com/pallandos/benchmarker/internal/monitor"
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

	if len(containers) == 0 {
		mainlogger.Warn("No containers found for stack: ", stackname)
		return
	}

	var monitorService *monitor.Service
	monitorService, err = monitor.NewService(cfg)
	if err != nil {
		mainlogger.Error("Failed to initialize monitor service: ", err)
		panic(err)
	}

	mainlogger.Info("Starting monitoring service...")
	monitorService.StartMonitoring(containers)

	// Create channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Create timer for monitoring duration
	timer := time.NewTimer(cfg.MonitorDuration)

	mainlogger.Infof("Monitoring started for %v. Press Ctrl+C to stop early...", cfg.MonitorDuration)

	// Wait for either interrupt signal or timer expiration
	select {
	case <-sigChan:
		mainlogger.Info("Interrupt signal received, shutting down...")
		timer.Stop()
	case <-timer.C:
		mainlogger.Info("Monitoring duration elapsed, shutting down...")
	}

	monitorService.Stop()
	mainlogger.Info("Shutdown complete")
}
