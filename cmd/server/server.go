package main

import (
	"github.com/sevlyar/go-daemon"
	log "github.com/sirupsen/logrus"
	"no_homomorphism/internal/app/server"
)

func main() {
	cntxt := &daemon.Context{
		PidFileName: "sample.pid",
		PidFilePerm: 0644,
		LogFileName: "sample.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[go-daemon sample]"},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()
	server.StartNew()
}
