package main

import (
	"github.com/VividCortex/godaemon"
	"no_homomorphism/internal/app/server"
)

func main() {
	godaemon.MakeDaemon(&godaemon.DaemonAttr{})
	server.StartNew()
}
