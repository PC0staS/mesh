package cmd

import (
	"fmt"

	"github.com/PC0staS/mesh/internal/client"
	"github.com/manifoldco/promptui"
)




func Remove() {
	// En lugar de leer config.json directamente,
	// pide la lista al daemon
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

	// Parsea servidores
	states, err := parseServerStates(response.Data)
	if err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		return
	}

	if len(states) == 0 {
		fmt.Println("No servers configured")
		return
	}

	// Crea lista de nombres para el select
	var serverNames []string
	for _, state := range states {
		serverNames = append(serverNames, state.Server.Name)
	}

	prompt := promptui.Select{
		Label: "Select Server to Remove",
		Items: serverNames,
	}
	index, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	// Envía request al daemon (con el índice)
	removeRequest := &client.Request{
		Command: "remove",
		Server:  float64(index),
	}

	removeResponse, err := client.SendRequest(removeRequest)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Is the daemon running? Try: mesh start")
		return
	}

	if !removeResponse.Success {
		fmt.Printf("❌ Error: %s\n", removeResponse.Message)
		return
	}

	fmt.Printf("✅ %s\n", removeResponse.Message)
}