package monitor

import "time"

type PingResult struct {
	Timestamp    time.Time     `json:"timestamp"`
	Success      bool          `json:"success"`
	ResponseTime time.Duration `json:"response_time"`
	Error        string        `json:"error,omitempty"`
}

type ServerState struct {
	Server         Server         `json:"server"`
	LastCheck      time.Time      `json:"last_check"`
	Status         bool           `json:"status"` // true = up, false = down
	ResponseTime   time.Duration  `json:"response_time"`
	LastError      string         `json:"last_error,omitempty"`
	UptimePercent  float64        `json:"uptime_percent"`
	AvgResponseTime time.Duration `json:"avg_response_time"`
	Results        []PingResult   `json:"results"` // Últimos N pings
}