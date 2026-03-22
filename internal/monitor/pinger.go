package monitor

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// PingServer hace un check a un servidor según su tipo
func PingServer(host string, checkType string, timeout time.Duration) PingResult {
	start := time.Now()

	switch checkType {
	case "ping":
		return pingTCP(host, timeout, start)
	case "http":
		return pingHTTP(host, timeout, start)
	case "tcp":
		return pingTCP(host, timeout, start)
	case "ssh":
		return pingSSH(host, timeout, start)
	default:
		return PingResult{
			Timestamp:    start,
			Success:      false,
			ResponseTime: 0,
			Error:        fmt.Sprintf("Unknown check type: %s", checkType),
		}
	}
}

// pingTCP intenta conectar por TCP
func pingTCP(host string, timeout time.Duration, start time.Time) PingResult {
	// Limpia la URL
	host = strings.TrimPrefix(host, "http://")
	host = strings.TrimPrefix(host, "https://")
	host = strings.TrimPrefix(host, "tcp://")

	// Si el host ya tiene puerto, úsalo. Si no, añade :80
	if !strings.Contains(host, ":") {
		host = host + ":80"
	}

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

// pingHTTP intenta hacer HTTP GET
func pingHTTP(host string, timeout time.Duration, start time.Time) PingResult {
	// Limpia la URL
	if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
		host = "http://" + host
	}

	// Crea cliente con timeout
	client := &net.Dialer{Timeout: timeout}
	conn, err := client.Dial("tcp", extractHostPort(host))
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

// pingSSH intenta conectar por SSH (puerto 22)
func pingSSH(host string, timeout time.Duration, start time.Time) PingResult {
	// Limpia la URL
	host = strings.TrimPrefix(host, "ssh://")

	// SSH por defecto en puerto 22
	if !strings.Contains(host, ":") {
		host = host + ":22"
	}

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

// extractHostPort extrae host:port de una URL
func extractHostPort(url string) string {
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")

	// Si no tiene puerto, añade :80
	if !strings.Contains(url, ":") {
		url = url + ":80"
	}

	return url
}