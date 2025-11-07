package filter

import (
	"channel_filter/config"
	"math"
	"time"
)

type ActionStats struct {
	EMA          float64
	Mean         float64
	StdDev       float64
	Count        int
	RecentValues []float64
	MaxSamples   int
	LastUpdate   time.Time
}

type UserBaseline struct {
	User              string
	ActionStats       map[string]*ActionStats
	LearningPhase     bool
	LearningEvents    int
	MinLearningEvents int
	WindowCounts      map[string]int
	EMAAalpha         float64
	StdDevMultiplier  float64
	DefaultThreshold  float64
	MinThreshold      float64
	MinSamples        int
}

func NewUserBaseline(user string, cfg *config.DetectionConfig) *UserBaseline {
	return &UserBaseline{
		User:              user,
		ActionStats:       make(map[string]*ActionStats),
		LearningPhase:     true,
		LearningEvents:    0,
		MinLearningEvents: cfg.MinLearningEvents,
		WindowCounts:      make(map[string]int),
		EMAAalpha:         cfg.EMAAalpha,
		StdDevMultiplier:  cfg.StdDevMultiplier,
		DefaultThreshold:  cfg.DefaultThreshold,
		MinThreshold:      cfg.MinThreshold,
		MinSamples:        cfg.MinSamplesForDetection,
	}
}

func (ub *UserBaseline) UpdateActionStats(actionType string, value float64, timestamp time.Time, maxSamples int) {
	stats, exists := ub.ActionStats[actionType]
	if !exists {
		stats = &ActionStats{
			EMA:          value,
			Mean:         value,
			StdDev:       0.0,
			Count:        1,
			RecentValues: make([]float64, 0),
			MaxSamples:   maxSamples,
			LastUpdate:   timestamp,
		}
		ub.ActionStats[actionType] = stats
	} else {
		stats.Count++
		stats.LastUpdate = timestamp
	}

	stats.EMA = ub.EMAAalpha*value + (1-ub.EMAAalpha)*stats.EMA
	stats.Mean = (stats.Mean*float64(stats.Count-1) + value) / float64(stats.Count)

	stats.RecentValues = append(stats.RecentValues, value)
	if len(stats.RecentValues) > stats.MaxSamples {
		stats.RecentValues = stats.RecentValues[1:]
	}

	if len(stats.RecentValues) >= 3 {
		recentMean := 0.0
		for _, v := range stats.RecentValues {
			recentMean += v
		}
		recentMean /= float64(len(stats.RecentValues))
		stats.StdDev = ub.calculateStdDev(stats.RecentValues, recentMean)
	} else {
		stats.StdDev = math.Abs(value-stats.Mean) * 0.5
	}
}

func (ub *UserBaseline) calculateStdDev(values []float64, mean float64) float64 {
	if len(values) < 2 {
		return 0.0
	}
	sumSquaredDiff := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquaredDiff += diff * diff
	}
	variance := sumSquaredDiff / float64(len(values))
	return math.Sqrt(variance)
}

func (ub *UserBaseline) IsAnomaly(actionType string, currentValue float64) bool {
	if ub.LearningPhase {
		return false
	}
	stats, exists := ub.ActionStats[actionType]
	if !exists {
		return false
	}
	if stats.Count < ub.MinSamples || stats.StdDev < 0.1 {
		return false
	}
	threshold := stats.Mean + (ub.StdDevMultiplier * stats.StdDev)
	emaThreshold := stats.EMA + (ub.StdDevMultiplier * stats.StdDev)
	return currentValue > threshold || currentValue > emaThreshold
}

func (ub *UserBaseline) GetAdaptiveThreshold(actionType string) float64 {
	if ub.LearningPhase {
		return ub.DefaultThreshold
	}
	stats, exists := ub.ActionStats[actionType]
	if !exists || stats.Count < ub.MinSamples {
		return ub.DefaultThreshold
	}
	threshold := stats.Mean + (ub.StdDevMultiplier * stats.StdDev)
	if threshold < ub.MinThreshold {
		threshold = ub.MinThreshold
	}
	return threshold
}

func (ub *UserBaseline) RecordEvent() {
	ub.LearningEvents++
	if ub.LearningEvents >= ub.MinLearningEvents {
		ub.LearningPhase = false
	}
}

func (ub *UserBaseline) UpdateWindowCount(actionType string, count int) {
	ub.WindowCounts[actionType] = count
}

func (ub *UserBaseline) GetWindowCount(actionType string) int {
	return ub.WindowCounts[actionType]
}
