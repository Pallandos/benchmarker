package stats

import (
	"context"
	"encoding/json"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type NetworkStats struct {
	ContainerID   string
	ContainerName string
	RxBytes       uint64
	TxBytes       uint64
	RxPackets     uint64
	TxPackets     uint64
	Timestamp     time.Time
}

type DockerMonitor struct {
	client *client.Client
}

func NewDockerMonitor() (*DockerMonitor, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &DockerMonitor{client: cli}, nil
}

func (dm *DockerMonitor) GetNetworkStats(ctx context.Context, containerID string) (*NetworkStats, error) {
	stats, err := dm.client.ContainerStats(ctx, containerID, false)
	if err != nil {
		return nil, err
	}
	defer stats.Body.Close()

	var containerStats types.StatsJSON
	if err := json.NewDecoder(stats.Body).Decode(&containerStats); err != nil {
		return nil, err
	}

	// Compute network statistics
	var rxBytes, txBytes, rxPackets, txPackets uint64
	for _, netStats := range containerStats.Networks {
		rxBytes += netStats.RxBytes
		txBytes += netStats.TxBytes
		rxPackets += netStats.RxPackets
		txPackets += netStats.TxPackets
	}

	return &NetworkStats{
		ContainerID:   containerID,
		ContainerName: containerStats.Name,
		RxBytes:       rxBytes,
		TxBytes:       txBytes,
		RxPackets:     rxPackets,
		TxPackets:     txPackets,
		Timestamp:     time.Now(),
	}, nil
}
