package cmd

import (
	"fmt"
	"strconv"

	"github.com/PC0staS/mesh/internal/client"
	"github.com/PC0staS/mesh/internal/monitor"
	"github.com/manifoldco/promptui"
)

func askName() string {
	prompt := promptui.Prompt{
		Label: "Server Name",
	}
	var name string
	for {
		n, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return ""
		}
		if n == "" {
			fmt.Println("Name cannot be empty")
			continue
		}
		name = n
		break
	}
	return name
}

func askHost() string {
	prompt := promptui.Prompt{
		Label: "Server Host (IP or URL with optional port)",
	}
	var host string
	for {
		h, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return ""
		}
		if h == "" {
			fmt.Println("Host cannot be empty")
			continue
		}
		host = h
		break
	}
	return host
}

func askType() string {
	prompt := promptui.Select{
		Label: "Server Type",
		Items: []string{"ping", "http", "tcp", "ssh"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "ping"
	}
	return result
}

func askInterval() int {
	prompt := promptui.Prompt{
		Label: "Interval (seconds)",
		Default: "60",
	}
	intervalStr, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return 0
	}
	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		fmt.Printf("Invalid interval %v\n, using default 60", err)
		return 60
	}
	return interval
}

func askTimeout() int {
	prompt := promptui.Prompt{
		Label: "Timeout (seconds)",
		Default: "5",
	}
	timeoutStr, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return 0
	}
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		fmt.Printf("Invalid timeout %v\n, using default 5", err)
		return 0
	}
	return timeout
}
func askEnabled() bool {
	prompt := promptui.Prompt{
		Label: "Enabled (Y/n)",
		Default: "Y",
	}
	enabledStr, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return false
	}
	if enabledStr == "" {
		enabledStr = "Y"
	}
	return enabledStr == "y" || enabledStr == "Y"
}
func askWebhook() string {
	prompt := promptui.Prompt{
		Label: "Webhook URL (leave empty for none)",
	}
	webhook, err := prompt.Run()
	if err != nil {
		return ""
	}
	
	return webhook
}

func Add(){
	server := monitor.Server{
		Name:          askName(),
		Host:          askHost(),
		Type:          askType(),
		Interval:      askInterval(),
		Timeout:       askTimeout(),
		Enabled:       askEnabled(),
		Webhook:       askWebhook(),
	}

	fmt.Printf("Adding server: %s (%s), type: %s, interval: %d, timeout: %d, enabled: %t\n", server.Name, server.Host, server.Type, server.Interval, server.Timeout, server.Enabled)

	// Envía request al daemon
	request := &client.Request{
		Command: "add",
		Server:  server, // Envía el struct completo
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