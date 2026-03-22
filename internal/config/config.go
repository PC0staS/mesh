package config

import "github.com/PC0staS/mesh/internal/monitor"

type Config struct {
	Servers []monitor.Server `json:"servers"`
}