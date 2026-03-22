package prompts

import (
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
)

func AskName() string {
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

func AskHost() string {
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

func AskType() string {
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

func AskInterval() int {
	for {
		prompt := promptui.Prompt{
			Label:   "Interval (seconds)",
			Default: "60",
		}
		intervalStr, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return 60
		}
		if intervalStr == "" {
			intervalStr = "60"
		}
		interval, err := strconv.Atoi(intervalStr)
		if err != nil {
			fmt.Println("❌ Invalid interval, please enter a number")
			continue
		}
		return interval
	}
}

func AskTimeout() int {
	for {
		prompt := promptui.Prompt{
			Label:   "Timeout (seconds)",
			Default: "5",
		}
		timeoutStr, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return 5
		}
		if timeoutStr == "" {
			timeoutStr = "5"
		}
		timeout, err := strconv.Atoi(timeoutStr)
		if err != nil {
			fmt.Println("❌ Invalid timeout, please enter a number")
			continue
		}
		return timeout
	}
}

func AskEnabled() bool {
	prompt := promptui.Prompt{
		Label:   "Enabled (Y/n)",
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

func AskWebhook() string {
	prompt := promptui.Prompt{
		Label: "Webhook URL (leave empty for none)",
	}
	webhook, err := prompt.Run()
	if err != nil {
		return ""
	}

	if webhook == "" {
		return ""
	}

	return webhook
}