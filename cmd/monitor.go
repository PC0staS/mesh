package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PC0staS/mesh/internal/client"
	"github.com/jroimartin/gocui"
)

var (
	serverStates []ServerState // Cache de estados
)

func Monitor() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	g.SelBgColor = gocui.ColorGreen
	g.SelFgColor = gocui.ColorBlack

	// Goroutine que actualiza cada 1 segundo
	go updateLoop(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Vista principal (tabla)
	if v, err := g.SetView("main", 0, 0, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "MESH - Server Status"
		v.Wrap = true
		v.Autoscroll = false
	}

	// Vista de status (abajo)
	if v, err := g.SetView("status", 0, maxY-3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Status"
		v.FgColor = gocui.ColorCyan
		fmt.Fprintln(v, "Press 'q' to quit")
	}

	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// updateLoop actualiza los datos cada X segundos
func updateLoop(g *gocui.Gui) {
	for {
		// Pide status al daemon
		request := &client.Request{
			Command: "status",
		}

		response, err := client.SendRequest(request)
		if err != nil {
			continue
		}

		if !response.Success {
			continue
		}

		// Parsea respuesta
		serverStates = parseStates(response.Data)

		// Renderiza
		g.Update(func(g *gocui.Gui) error {
			v, _ := g.View("main")
			v.Clear()

			// Headers
			fmt.Fprintf(v, "%-20s %-25s %-10s %-8s %-10s\n", "Name", "Host", "Type", "Status", "Uptime")
			fmt.Fprintln(v, strings.Repeat("-", 75))

			// Servidores
			for _, state := range serverStates {
				status := "DOWN"
				if state.Status {
					status = "UP"
				}
				fmt.Fprintf(v, "%-20s %-25s %-10s %-8s %.1f%%\n",
					state.Server.Name, state.Server.Host, state.Server.Type, status, state.UptimePercent)
			}

			return nil
		})

		time.Sleep(1 * time.Second)
	}
}

func parseStates(data interface{}) []ServerState {
	rawStates, ok := data.([]interface{})
	if !ok {
		return []ServerState{}
	}

	var states []ServerState
	for _, raw := range rawStates {
		stateMap, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		var server ServerStruct
		if serverMap, ok := stateMap["server"].(map[string]interface{}); ok {
			server = ServerStruct{
				Name:     getString(serverMap, "name"),
				Host:     getString(serverMap, "host"),
				Type:     getString(serverMap, "type"),
				Interval: getInt(serverMap, "interval"),
			}
		}

		state := ServerState{
			Server:        server,
			Status:        getBool(stateMap, "status"),
			UptimePercent: getFloat(stateMap, "uptime_percent"),
		}

		states = append(states, state)
	}

	return states
}