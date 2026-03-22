package monitor

import (
	"net"
	"time"
)

// PingServer hace un ping a un servidor (usando TCP)
func PingServer(host string, timeout time.Duration) PingResult {
	start := time.Now()
	
	// Intenta conectar por TCP al puerto 80 (HTTP)
	conn, err := net.DialTimeout("tcp", host+":80", timeout)
	if err != nil {
		return PingResult{
			Timestamp:    start,
			Success:      false,
			ResponseTime: 0,
			Error:        err.Error(),
		}
	}
	defer conn.Close()

	elapsed := time.Since(start)
	return PingResult{
		Timestamp:    start,
		Success:      true,
		ResponseTime: elapsed,
	}
}