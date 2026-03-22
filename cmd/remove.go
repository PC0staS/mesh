package cmd

import (
	"fmt"

	"github.com/PC0staS/mesh/internal/client"
	"github.com/PC0staS/mesh/internal/config"
	"github.com/manifoldco/promptui"
)




func Remove() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}
	if len(cfg.Servers) == 0 {
		fmt.Println("No servers configured")
		return
	}

	// Select server to remove
	prompt := promptui.Select{
		Label: "Select Server to Remove",
		Items: cfg.Servers,
		Templates: &promptui.SelectTemplates{
			Label: "{{ . }}",
			Active:   "▸ {{ .Name | cyan }} ({{ .Host }}) [{{ .Type }}]",
			Inactive: "  {{ .Name }} ({{ .Host }}) [{{ .Type }}]",
			Selected: "Removed {{ .Name }} ({{ .Host }}) [{{ .Type }}]",
		},
	}
	index, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	// Envía request al daemon (con el índice)
	request := &client.Request{
		Command: "remove",
		Server:  float64(index), // Envía el índice como float64
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