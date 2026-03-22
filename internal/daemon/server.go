package daemon

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/PC0staS/mesh/internal/client"
	"github.com/PC0staS/mesh/internal/config"
	"github.com/PC0staS/mesh/internal/monitor"
)

const SocketPath = "/tmp/mesh.sock"

// StartDaemon inicia el servidor que escucha en el socket
func StartDaemon() {
	// Carga config
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	if len(cfg.Servers) == 0 {
		fmt.Println("No servers configured. Add some with: mesh add")
		return
	}

	// Inicia pinging
	StartPinging(cfg)
	// 1. Elimina socket anterior si existe
	os.Remove(SocketPath)

	// 2. Crea listener en el socket
	listener, err := net.Listen("unix", SocketPath)
	if err != nil {
		fmt.Printf("Error starting daemon: %v\n", err)
		return
	}
	defer listener.Close()
	defer os.Remove(SocketPath)

	fmt.Printf("✅ Daemon started, listening on %s\n", SocketPath)

	// 3. Loop infinito: acepta conexiones
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		// 4. Maneja cada conexión en goroutine
		go handleConnection(conn)
	}
}

func StopDaemon() {
	// Para detener el daemon, simplemente eliminamos el socket
	err := os.Remove(SocketPath)
	if err != nil {
		fmt.Printf("Error stopping daemon: %v\n", err)
		return
	}
	fmt.Println("Daemon stopped")
}

// handleConnection procesa un request del cliente
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Lee request
	var request client.Request
	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&request)
	if err != nil {
		fmt.Printf("Error decoding request: %v\n", err)
		return
	}

	// Procesa según comando
	var response *client.Response

	switch request.Command {
	case "status":
		response = handleStatus()
	case "add":
		response = handleAdd(request)
	case "remove":
		response = handleRemove(request)
	default:
		response = &client.Response{
			Success: false,
			Message: "Unknown command",
		}
	}

	// Envía respuesta
	encoder := json.NewEncoder(conn)
	encoder.Encode(response)
}

// handleStatus devuelve el estado actual
func handleStatus() *client.Response {
	states := daemonState.GetAllStates()

	return &client.Response{
		Success: true,
		Data:    states,
	}
}

// handleAdd añade un servidor
func handleAdd(request client.Request) *client.Response {
	// Convierte request.Server a monitor.Server
	serverMap, ok := request.Server.(map[string]interface{})
	if !ok {
		return &client.Response{
			Success: false,
			Message: "Invalid server data",
		}
	}

	server := monitor.Server{
		Name:     getString(serverMap, "name"),
		Host:     getString(serverMap, "host"),
		Type:     getString(serverMap, "type"),
		Interval: getInt(serverMap, "interval"),
		Timeout:  getInt(serverMap, "timeout"),
		Enabled:  getBool(serverMap, "enabled"),
	}

	// Carga config actual
	cfg, err := config.LoadConfig()
	if err != nil {
		return &client.Response{
			Success: false,
			Message: fmt.Sprintf("Error loading config: %v", err),
		}
	}

	// Añade servidor
	cfg.Servers = append(cfg.Servers, server)

	// Guarda config
	err = config.SaveConfig(cfg)
	if err != nil {
		return &client.Response{
			Success: false,
			Message: fmt.Sprintf("Error saving config: %v", err),
		}
	}

	return &client.Response{
		Success: true,
		Message: fmt.Sprintf("Server '%s' added successfully", server.Name),
	}
}

// Funciones helper (copialas de cmd/list.go)
func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getInt(m map[string]interface{}, key string) int {
	if v, ok := m[key].(float64); ok {
		return int(v)
	}
	return 0
}

func getBool(m map[string]interface{}, key string) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
}

// handleRemove elimina un servidor
func handleRemove(request client.Request) *client.Response {
	// Convierte request.Server a int (índice)

	// Esto es un poco feo, pero funciona
	index, ok := request.Server.(float64)
	if !ok {
		return &client.Response{
			Success: false,
			Message: "Invalid index",
		}
	}

	idx := int(index)

	// Carga config
	cfg, err := config.LoadConfig()
	if err != nil {
		return &client.Response{
			Success: false,
			Message: fmt.Sprintf("Error loading config: %v", err),
		}
	}

	// Valida índice
	if idx < 0 || idx >= len(cfg.Servers) {
		return &client.Response{
			Success: false,
			Message: "Invalid server index",
		}
	}

	// Guarda nombre para mensaje
	serverName := cfg.Servers[idx].Name

	// Elimina servidor
	cfg.Servers = append(cfg.Servers[:idx], cfg.Servers[idx+1:]...)

	// Guarda config
	err = config.SaveConfig(cfg)
	if err != nil {
		return &client.Response{
			Success: false,
			Message: fmt.Sprintf("Error saving config: %v", err),
		}
	}

	return &client.Response{
		Success: true,
		Message: fmt.Sprintf("Server '%s' removed successfully", serverName),
	}
}


var daemonState *DaemonState

// StartPinging inicia goroutines para hacer ping a cada servidor
func StartPinging(cfg *config.Config) {
	daemonState = NewDaemonState(100) // Máx 100 pings por servidor

	// Inicializa estado para cada servidor
	for _, server := range cfg.Servers {
		if !server.Enabled {
			continue
		}

		// Crea ServerState inicial
		serverState := &monitor.ServerState{
			Server:      server,
			Status:      false,
			Results:     []monitor.PingResult{},
			UptimePercent: 100.0,
		}

		daemonState.Servers[server.Name] = serverState

		// Lanza goroutine para hacer pings
		go pingLoop(server)
	}

	fmt.Println("✅ Pinging started for all servers")
}

// pingLoop hace pings continuamente a un servidor
func pingLoop(server monitor.Server) {
	ticker := time.NewTicker(time.Duration(server.Interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Hace ping
		result := monitor.PingServer(server.Host, time.Duration(server.Timeout)*time.Second)

		// Añade resultado al estado
		daemonState.AddResult(server.Name, result)

		// Log (opcional, para debugging)
		if result.Success {
			fmt.Printf("[%s] ✅ %s (%.2fms)\n", server.Name, server.Host, float64(result.ResponseTime.Milliseconds()))
		} else {
			fmt.Printf("[%s] ❌ %s (Error: %s)\n", server.Name, server.Host, result.Error)
		}
	}
}