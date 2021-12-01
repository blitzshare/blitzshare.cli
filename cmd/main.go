package main

import (
	"flag"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/services"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/services/str"
	"os"
	"os/signal"
	"syscall"

	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/config"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/dependencies"
	log "github.com/sirupsen/logrus"
)

func initLog() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	initLog()
	peer := flag.Bool("peer", false, "Connect to p2p peer")
	flag.Parse()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config %v\n", err)
	}
	deps, err := dependencies.NewDependencies(cfg)
	if err != nil {
		log.Fatalf("failed to load dependencies %v\n", err)
	}
	if *peer {
		app.StartPeer(deps)
	} else {
		log.Println("Enter OTP:")
		line := services.ReadStdInLine()
		otp := str.SanatizeStr(*line)
		app.ConnectToPeerOTP(deps, &otp)
	}
	services.PrintLogo()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)
	select {
	case <-stop:
		deps.P2p.Close()
		os.Exit(0)
	}
}
