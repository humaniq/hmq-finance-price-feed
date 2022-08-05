package config

import "time"

const DefaultFilterPercentThreshold float64 = 1

type StorageConsumer struct {
	Name        string          `yaml:"name,omitempty"`
	Datastore   GoogleDataStore `yaml:"google_datastore,omitempty"`
	Thresholds  Thresholds      `yaml:"thresholds,omitempty"`
	TimeDiffStr string          `yaml:"time_diff"`
}

func (sc *StorageConsumer) TimeDelta() time.Duration {
	duration, err := time.ParseDuration(sc.TimeDiffStr)
	if err != nil {
		return 0
	}
	return duration
}

type Thresholds struct {
	Default float64            `yaml:"default,omitempty"`
	Symbols map[string]float64 `yaml:"symbols"`
}
