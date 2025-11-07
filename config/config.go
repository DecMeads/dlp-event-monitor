package config

import "time"

type Config struct {
	Detection DetectionConfig
	Window    WindowConfig
}

type DetectionConfig struct {
	MinLearningEvents      int
	StdDevMultiplier       float64
	EMAAalpha              float64
	MaxSamples             int
	MinSamplesForDetection int
	DefaultThreshold       float64
	MinThreshold           float64
}

type WindowConfig struct {
	Duration time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		Detection: DetectionConfig{
			MinLearningEvents:      30,
			StdDevMultiplier:       3.0,
			EMAAalpha:              0.1,
			MaxSamples:             20,
			MinSamplesForDetection: 5,
			DefaultThreshold:       10.0,
			MinThreshold:           2.0,
		},
		Window: WindowConfig{
			Duration: 5 * time.Minute,
		},
	}
}
