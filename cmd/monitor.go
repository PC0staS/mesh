package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
	"sort"

	"github.com/PC0staS/mesh/internal/client"
	"github.com/PC0staS/mesh/internal/monitor"
	"github.com/jroimartin/gocui"
)

var (
	serverStates  []monitor.ServerState
	selectedIndex int = 0
)

func Monitor() {
	checkRoot() // Solo root puede monitorear (porque muestra status)
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

	go updateLoop(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Vista lista (izquierda)
	if v, err := g.SetView("list", 0, 0, maxX/2-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Servers"
		v.Wrap = false
		v.Autoscroll = false
	}

	// Vista detalles (derecha)
	if v, err := g.SetView("details", maxX/2, 0, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Details"
		v.Wrap = true
		v.Autoscroll = true
	}

	// Vista status (abajo)
	if v, err := g.SetView("status", 0, maxY-3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Status"
		v.FgColor = gocui.ColorCyan
		fmt.Fprintln(v, "↑↓: Navigate | q: Quit")
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
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, selectUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, selectDown); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func selectUp(g *gocui.Gui, v *gocui.View) error {
	if selectedIndex > 0 {
		selectedIndex--
	}
	return nil
}

func selectDown(g *gocui.Gui, v *gocui.View) error {
	if selectedIndex < len(serverStates)-1 {
		selectedIndex++
	}
	return nil
}

func updateLoop(g *gocui.Gui) {
	for {
		request := &client.Request{
			Command: "status",
		}

		response, err := client.SendRequest(request)
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		if !response.Success {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		serverStates = parseStates(response.Data)
		sort.Slice(serverStates, func(i, j int) bool {
			return serverStates[i].Server.Name < serverStates[j].Server.Name
		})

		g.Update(func(g *gocui.Gui) error {
			// Renderiza lista
			listView, _ := g.View("list")
			listView.Clear()

			fmt.Fprintf(listView, "%-20s %-8s %-10s\n", "Name", "Status", "Uptime")
			fmt.Fprintln(listView, strings.Repeat("-", 40))

			for i, state := range serverStates {
				status := "🔴"
				if state.Status {
					status = "🟢"
				}

				selected := " "
				if i == selectedIndex {
					selected = ">"
				}

				fmt.Fprintf(listView, "%s %-18s %-8s %.1f%%\n",
					selected, state.Server.Name, status, state.UptimePercent)
			}

			// Renderiza detalles
			detailsView, _ := g.View("details")
			detailsView.Clear()

			if selectedIndex < len(serverStates) {
				selected := serverStates[selectedIndex]
				fmt.Fprintf(detailsView, "Name: %s\n", selected.Server.Name)
				fmt.Fprintf(detailsView, "Host: %s\n", selected.Server.Host)
				fmt.Fprintf(detailsView, "Type: %s\n", selected.Server.Type)
				fmt.Fprintf(detailsView, "Status: %v\n", selected.Status)
				fmt.Fprintf(detailsView, "Uptime: %.1f%%\n", selected.UptimePercent)
				fmt.Fprintln(detailsView, "")
				fmt.Fprintln(detailsView, "Recent checks:")
				fmt.Fprintln(detailsView, strings.Repeat("-", 30))

				// Muestra últimos 10 resultados
				start := len(selected.Results) - 10
				if start < 0 {
					start = 0
				}

				for _, result := range selected.Results[start:] {
					status := "❌"
					if result.Success {
						status = "✅"
					}

					fmt.Fprintf(detailsView, "%s %.2fms\n", status, float64(result.ResponseTime.Milliseconds()))
				}
			}

			return nil
		})

		time.Sleep(500 * time.Millisecond)
	}
}

func parseStates(data interface{}) []monitor.ServerState {
	jsonBytes, _ := json.Marshal(data)
	var states []monitor.ServerState
	json.Unmarshal(jsonBytes, &states)
	return states
}