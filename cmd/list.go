package cmd

import (
	"fmt"

	"github.com/PC0staS/mesh/internal/config"
)

func List() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}
	if len(cfg.Servers) == 0 {
		fmt.Println("No servers configured")
		return
	}
	fmt.Println("Configured Servers:")
	// Headers
	fmt.Printf("%-20s %-30s %-10s %-10s %-10s %-10s\n", "Name", "Host", "Type", "Interval", "Timeout", "Enabled")
	for _, server := range cfg.Servers {
		fmt.Printf("%-20s %-30s %-10s %-10d %-10d %-10t\n", server.Name, server.Host, server.Type, server.Interval, server.Timeout, server.Enabled)
	}
}