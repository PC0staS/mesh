package cmd

import (
	"fmt"
	"strconv"

	"github.com/PC0staS/mesh/internal/config"
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

func Add(){
	var server monitor.Server

	server = monitor.Server{
	Name: askName(),
	Host: askHost(),
	Type: askType(),
	Interval: askInterval(),
	Timeout: askTimeout(),
	Enabled: askEnabled(),
	}

	fmt.Printf("Adding server: %s (%s), type: %s, interval: %d, timeout: %d, enabled: %t\n", server.Name, server.Host, server.Type, server.Interval, server.Timeout, server.Enabled)

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	cfg.Servers = append(cfg.Servers, server)

	err = config.SaveConfig(cfg)
	if err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		return
	}
	fmt.Println("Server added successfully")
}