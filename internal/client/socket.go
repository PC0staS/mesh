package client

import (
	"encoding/json"
	"fmt"
	"net"
)

const SocketPath = "/tmp/mesh.sock"

// SendRequest envía un request al daemon y espera respuesta
func SendRequest(request *Request) (*Response, error) {
	// 1. Conecta al socket
	conn, err := net.Dial("unix", SocketPath)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to daemon: %v", err)
	}
	defer conn.Close()

	// 2. Envía request (JSON)
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(request)
	if err != nil {
		return nil, fmt.Errorf("cannot send request: %v", err)
	}

	// 3. Lee respuesta (JSON)
	var response Response
	decoder := json.NewDecoder(conn)
	err = decoder.Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("cannot read response: %v", err)
	}

	return &response, nil
}