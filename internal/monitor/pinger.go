package monitor

import (
	"net"
	"strings"
	"time"
)

// PingServer hace un ping a un servidor (usando TCP)
func PingServer(host string, timeout time.Duration) PingResult {
	start := time.Now()
	
	// Limpia la URL (quita http://, https://, etc.)
	host = strings.TrimPrefix(host, "http://")
	host = strings.TrimPrefix(host, "https://")
	host = strings.TrimPrefix(host, "tcp://")
	
	// Si el host ya tiene puerto, úsalo. Si no, añade :80
	if !strings.Contains(host, ":") {
		host = host + ":80"
	}
	
	// Intenta conectar por TCP
	conn, err := net.DialTimeout("tcp", host, timeout)
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