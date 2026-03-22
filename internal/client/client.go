package client

type Request struct {
	Command string      `json:"command"` // "status", "add", "remove"
	Server  interface{} `json:"server,omitempty"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}