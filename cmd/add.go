package cmd

import (
	"fmt"

	"github.com/PC0staS/mesh/internal/client"
	"github.com/PC0staS/mesh/internal/monitor"
	"github.com/PC0staS/mesh/internal/prompts"
)

func Add() {
	var server monitor.Server

	server = monitor.Server{
		Name:     prompts.AskName(),
		Host:     prompts.AskHost(),
		Type:     prompts.AskType(),
		Interval: prompts.AskInterval(),
		Timeout:  prompts.AskTimeout(),
		Enabled:  prompts.AskEnabled(),
		Webhook:  prompts.AskWebhook(),
	}

	fmt.Printf("Adding server: %s (%s), type: %s, interval: %d, timeout: %d, enabled: %t\n",
		server.Name, server.Host, server.Type, server.Interval, server.Timeout, server.Enabled)

	request := &client.Request{
		Command: "add",
		Server:  server,
	}

	response, err := client.SendRequest(request)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Is the daemon running? Try: mesh start")
		return
	}

	if !response.Success {
		fmt.Printf("❌ Error: %s\n", response.Message)
		return
	}

	fmt.Printf("✅ %s\n", response.Message)
}