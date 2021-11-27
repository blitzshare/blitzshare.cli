package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/config"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/dependencies"
	"github.com/libp2p/go-libp2p-core/host"
	log "github.com/sirupsen/logrus"
)

func initLog() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	// initLog()
	dest := flag.String("d", "", "Destination multiaddr string")
	pass := flag.String("p", "", "One time pass id of connection peer")
	flag.Parse()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config %v\n", err)
	}
	deps, err := dependencies.NewDependencies(cfg)
	if err != nil {
		log.Fatalf("failed to load dependencies %v\n", err)
	}

	var host host.Host
	if *dest != "" {
		host = app.ConnectToPeerAddress(deps, dest)
	} else if *pass != "" {
		host = app.ConnectToPeerPass(deps, pass)
	} else {
		host = app.StartPeer(deps)
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)
	select {
	case <-stop:
		host.Close()
		os.Exit(0)
	}
}
