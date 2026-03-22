package cmd

import "github.com/PC0staS/mesh/internal/daemon"

func Stop() {
	daemon.StopDaemon()
}