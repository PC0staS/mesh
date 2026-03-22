package cmd

import (
	"fmt"

	"github.com/PC0staS/mesh/internal/client"
)

func Reload() {
	checkRoot() // Solo root puede recargar (porque muestra status)
	request := &client.Request{
		Command: "reload",
	}

	response, err := client.SendRequest(request)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if !response.Success {
		fmt.Printf("❌ Error: %s\n", response.Message)
		return
	}

	fmt.Printf("✅ %s\n", response.Message)
}