package cmd

import (
	"fmt"

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
	serverToRemove := cfg.Servers[index]

	// Remove server from config
	cfg.Servers = append(cfg.Servers[:index], cfg.Servers[index+1:]...)
	err = config.SaveConfig(cfg)
	if err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		return
	}
	fmt.Printf("Server removed: %s (%s), type: %s\n", serverToRemove.Name, serverToRemove.Host, serverToRemove.Type)
}