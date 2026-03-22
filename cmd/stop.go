package cmd

import "github.com/PC0staS/mesh/internal/daemon"

func Stop() {
	checkRoot()
	daemon.StopDaemon()
}