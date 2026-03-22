package daemon

import (
	"sync"
	"time"

	"github.com/PC0staS/mesh/internal/monitor"
)

type DaemonState struct {
	mu       sync.RWMutex                    // Para acceso thread-safe
	Servers  map[string]*monitor.ServerState // Indexado por nombre
	maxResults int                           // Máximo de histórico por servidor
}

func NewDaemonState(maxResults int) *DaemonState {
	return &DaemonState{
		Servers:    make(map[string]*monitor.ServerState),
		maxResults: maxResults,
	}
}

// AddResult añade un ping result al histórico de un servidor
func (ds *DaemonState) AddResult(serverName string, result monitor.PingResult) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	state, exists := ds.Servers[serverName]
	if !exists {
		return // Servidor no existe
	}

	// Añade resultado
	state.Results = append(state.Results, result)

	// Limita histórico
	if len(state.Results) > ds.maxResults {
		state.Results = state.Results[len(state.Results)-ds.maxResults:]
	}

	// Actualiza estado actual
	state.LastCheck = result.Timestamp
	state.Status = result.Success
	state.ResponseTime = result.ResponseTime
	if result.Error != "" {
		state.LastError = result.Error
	}

	// Calcula uptime %
	state.UptimePercent = calculateUptime(state.Results)
	state.AvgResponseTime = calculateAvgResponseTime(state.Results)
}

// GetServerState devuelve el estado actual de un servidor
func (ds *DaemonState) GetServerState(serverName string) *monitor.ServerState {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	return ds.Servers[serverName]
}

// GetAllStates devuelve todos los servidores
func (ds *DaemonState) GetAllStates() []*monitor.ServerState {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	states := make([]*monitor.ServerState, 0)
	for _, state := range ds.Servers {
		states = append(states, state)
	}
	return states
}

func calculateUptime(results []monitor.PingResult) float64 {
	if len(results) == 0 {
		return 100.0
	}

	ups := 0
	for _, r := range results {
		if r.Success {
			ups++
		}
	}

	return (float64(ups) / float64(len(results))) * 100
}

func calculateAvgResponseTime(results []monitor.PingResult) time.Duration {
	if len(results) == 0 {
		return 0
	}

	var total time.Duration
	count := 0

	for _, r := range results {
		if r.Success {
			total += r.ResponseTime
			count++
		}
	}

	if count == 0 {
		return 0
	}

	return total / time.Duration(count)
}