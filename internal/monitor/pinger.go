package monitor

import (
	"fmt"
	"net"
	"net/http"
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
	host = strings.TrimPrefix(host, "http://")
	host = strings.TrimPrefix(host, "https://")
	host = strings.TrimPrefix(host, "tcp://")

	// Quita ruta si existe
	if idx := strings.Index(host, "/"); idx != -1 {
		host = host[:idx]
	}

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

// pingHTTP intenta hacer conexión HTTP
func pingHTTP(host string, timeout time.Duration, start time.Time) PingResult {
	original := host // Guarda original para detectar https
	
	// Limpia esquemas
	isHTTPS := strings.Contains(original, "https://")
	host = strings.TrimPrefix(host, "http://")
	host = strings.TrimPrefix(host, "https://")
	host = strings.TrimPrefix(host, "/")
	
	// Quita ruta si existe
	if idx := strings.Index(host, "/"); idx != -1 {
		host = host[:idx]
	}
	
	// Puerto por defecto según protocolo
	if !strings.Contains(host, ":") {
		if isHTTPS {
			host = host + ":443"
		} else {
			host = host + ":80"
		}
	}

	// Construye URL
	scheme := "http://"
	if isHTTPS {
		scheme = "https://"
	}
	url := scheme + host

	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return PingResult{
			Timestamp:    start,
			Success:      false,
			ResponseTime: 0,
			Error:        err.Error(),
		}
	}
	defer resp.Body.Close()

	// Check status 200-299
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return PingResult{
			Timestamp:    start,
			Success:      false,
			ResponseTime: time.Since(start),
			Error:        fmt.Sprintf("HTTP %d", resp.StatusCode),
		}
	}

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