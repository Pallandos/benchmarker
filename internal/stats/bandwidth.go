package stats

import (
	"time"
)

type BandwidthMetrics struct {
	ContainerID     string
	ContainerName   string
	RxBytesPerSec   float64
	TxBytesPerSec   float64
	RxPacketsPerSec float64
	TxPacketsPerSec float64
	Period          time.Duration
}

type BandwidthCalculator struct {
	previousStats map[string]*NetworkStats
}

func NewBandwidthCalculator() *BandwidthCalculator {
	return &BandwidthCalculator{
		previousStats: make(map[string]*NetworkStats),
	}
}

func (bc *BandwidthCalculator) Calculate(current *NetworkStats) *BandwidthMetrics {
	previous, exists := bc.previousStats[current.ContainerID]

	if !exists {
		bc.previousStats[current.ContainerID] = current
		return &BandwidthMetrics{
			ContainerID:   current.ContainerID,
			ContainerName: current.ContainerName,
			Period:        0,
		}
	}

	timeDiff := current.Timestamp.Sub(previous.Timestamp).Seconds()
	if timeDiff <= 0 {
		return nil
	}

	metrics := &BandwidthMetrics{
		ContainerID:     current.ContainerID,
		ContainerName:   current.ContainerName,
		RxBytesPerSec:   float64(current.RxBytes-previous.RxBytes) / timeDiff,
		TxBytesPerSec:   float64(current.TxBytes-previous.TxBytes) / timeDiff,
		RxPacketsPerSec: float64(current.RxPackets-previous.RxPackets) / timeDiff,
		TxPacketsPerSec: float64(current.TxPackets-previous.TxPackets) / timeDiff,
		Period:          current.Timestamp.Sub(previous.Timestamp),
	}

	bc.previousStats[current.ContainerID] = current
	return metrics
}
