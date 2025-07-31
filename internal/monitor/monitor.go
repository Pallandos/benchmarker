package monitor

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/pallandos/benchmarker/internal/containers"
	"github.com/pallandos/benchmarker/internal/stats"
	"github.com/pallandos/benchmarker/internal/utils/config"
	"github.com/pallandos/benchmarker/internal/utils/logger"
)

type Service struct {
	dockerMonitor *stats.DockerMonitor
	bandwidthCalc *stats.BandwidthCalculator
	config        *config.AppConfig
	logger        *logrus.Logger
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

func NewService(cfg *config.AppConfig) (*Service, error) {
	dockerMonitor, err := stats.NewDockerMonitor()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	monitorLogger, err := logger.InitLogger("monitor.log", cfg.LogPath)

	if err != nil {
		cancel()
		return nil, err
	}

	return &Service{
		dockerMonitor: dockerMonitor,
		bandwidthCalc: stats.NewBandwidthCalculator(),
		config:        cfg,
		logger:        monitorLogger,
		ctx:           ctx,
		cancel:        cancel,
	}, nil
}

func (s *Service) StartMonitoring(containerInfos []containers.ContainerInfo) {
	for _, containerInfo := range containerInfos {
		s.wg.Add(1)
		go s.monitorContainer(containerInfo)
	}
}

func (s *Service) Stop() {
	s.cancel()
	s.wg.Wait()
}

func (s *Service) monitorContainer(containerInfo containers.ContainerInfo) {
	defer s.wg.Done()

	ticker := time.NewTicker(s.config.MonitorInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			networkStats, err := s.dockerMonitor.GetNetworkStats(s.ctx, containerInfo.ID)
			if err != nil {
				s.logger.WithFields(logrus.Fields{
					"container_id":   containerInfo.ID,
					"container_name": containerInfo.Name,
					"error":          err,
				}).Error("Failed to get network stats")
				continue
			}

			bandwidthMetrics := s.bandwidthCalc.Calculate(networkStats)
			if bandwidthMetrics != nil && bandwidthMetrics.Period > 0 {
				s.logBandwidthMetrics(bandwidthMetrics, containerInfo.Name)
			}
		}
	}
}

func (s *Service) logBandwidthMetrics(metrics *stats.BandwidthMetrics, containerName string) {
	s.logger.WithFields(logrus.Fields{
		"container_id":       metrics.ContainerID,
		"container_name":     containerName,
		"rx_bytes_per_sec":   metrics.RxBytesPerSec,
		"tx_bytes_per_sec":   metrics.TxBytesPerSec,
		"rx_packets_per_sec": metrics.RxPacketsPerSec,
		"tx_packets_per_sec": metrics.TxPacketsPerSec,
		"period_ms":          metrics.Period.Milliseconds(),
	}).Info("Bandwidth metrics")
}
