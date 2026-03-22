package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PC0staS/mesh/internal/client"
	"github.com/PC0staS/mesh/internal/monitor"
)

func List() {
	checkRoot() // Solo root puede listar (porque muestra status)
	// Primero, recarga la config en el daemon
	reloadRequest := &client.Request{
		Command: "reload",
	}
	client.SendRequest(reloadRequest)
	time.Sleep(100 * time.Millisecond)

	request := &client.Request{
		Command: "status",
	}

	response, err := client.SendRequest(request)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Is the daemon running? Try: mesh start")
		return
	}

	if !response.Success {
		fmt.Printf("Error: %s\n", response.Message)
		return
	}

	// Parsea directamente a []monitor.ServerState
	jsonBytes, _ := json.Marshal(response.Data)
	var states []monitor.ServerState
	json.Unmarshal(jsonBytes, &states)

	if len(states) == 0 {
		fmt.Println("No servers configured")
		return
	}

	// Headers
	fmt.Printf("%-20s %-25s %-10s %-10s %-8s %-10s\n", "Name", "Host", "Type", "Interval", "Status", "Uptime")
	fmt.Println(strings.Repeat("-", 85))

	// Servidores
	for _, state := range states {
		status := "❌"
		if state.Status {
			status = "✅"
		}
		fmt.Printf("%-20s %-25s %-10s %-10d %-8s %.1f%%\n",
			state.Server.Name, state.Server.Host, state.Server.Type, state.Server.Interval, status, state.UptimePercent)
	}
}

type ServerState struct {
	Server        ServerStruct `json:"server"`
	LastCheck     string       `json:"last_check"`
	Status        bool         `json:"status"`
	ResponseTime  string       `json:"response_time"`
	UptimePercent float64      `json:"uptime_percent"`
	AvgResponseTime string      `json:"avg_response_time"`
}

type ServerStruct struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Type     string `json:"type"`
	Interval int    `json:"interval"`
	Timeout  int    `json:"timeout"`
	Enabled  bool   `json:"enabled"`
}

func parseServerStates(data interface{}) ([]ServerState, error) {
	rawStates, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data type")
	}

	var states []ServerState
	for _, raw := range rawStates {
		stateMap, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		// Parsea nested Server
		var server ServerStruct
		if serverMap, ok := stateMap["server"].(map[string]interface{}); ok {
			server = ServerStruct{
				Name:     getString(serverMap, "name"),
				Host:     getString(serverMap, "host"),
				Type:     getString(serverMap, "type"),
				Interval: getInt(serverMap, "interval"),
				Timeout:  getInt(serverMap, "timeout"),
				Enabled:  getBool(serverMap, "enabled"),
			}
		}

		state := ServerState{
			Server:        server,
			Status:        getBool(stateMap, "status"),
			UptimePercent: getFloat(stateMap, "uptime_percent"),
		}

		states = append(states, state)
	}

	return states, nil
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getInt(m map[string]interface{}, key string) int {
	if v, ok := m[key].(float64); ok {
		return int(v)
	}
	return 0
}

func getBool(m map[string]interface{}, key string) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
}

func getFloat(m map[string]interface{}, key string) float64 {
	if v, ok := m[key].(float64); ok {
		return v
	}
	return 0.0
}