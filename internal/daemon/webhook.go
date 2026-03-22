package daemon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/PC0staS/mesh/internal/monitor"
)

// SendWebhook envía notificación al webhook
func SendWebhook(server monitor.Server, state *monitor.ServerState) {
	if server.Webhook == "" {
		return
	}

	// Solo JSON, siempre
	payload := map[string]interface{}{
		"server_name":    server.Name,
		"host":           server.Host,
		"status":         state.Status,
		"timestamp":      time.Now(),
		"response_time":  state.ResponseTime.Milliseconds(),
		"uptime_percent": state.UptimePercent,
		"avg_response":   state.AvgResponseTime.Milliseconds(),
	}
	
	body, _ := json.Marshal(payload)

	go func() {
		client := &http.Client{Timeout: 5 * time.Second}
		req, err := http.NewRequest("POST", server.Webhook, bytes.NewBuffer(body))
		if err != nil {
			fmt.Printf("Error creating webhook request: %v\n", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error sending webhook: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			fmt.Printf("✅ Webhook sent for %s\n", server.Name)
		} else {
			fmt.Printf("❌ Webhook failed for %s (status %d)\n", server.Name, resp.StatusCode)
		}
	}()
}