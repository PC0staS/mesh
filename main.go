package main

import (
	"fmt"
	"os"

	"github.com/PC0staS/mesh/cmd"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "help" {
		fmt.Print("MESH - Monitor Each Server Health\n")
		fmt.Println("Usage: mesh <command>")
		fmt.Println("Commands: add, list, remove, start, stop, monitor")
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		cmd.Add()
	case "list":
		cmd.List()
	case "remove":
		cmd.Remove()
	case "start":
		cmd.Start()
	case "stop":
		cmd.Stop()
	case "monitor":
		cmd.Monitor()
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}