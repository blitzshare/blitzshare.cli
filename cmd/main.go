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
	"github.com/libp2p/go-libp2p-core/host"
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
	var host *host.Host
	if *peer {
		host = app.StartPeer(deps)
	}else {
		log.Println("Enter OTP:")
		line := services.ReadStdInLine()
		otp := str.SanatizeStr(*line)
		host = app.ConnectToPeerPass(deps, &otp)
	}
	services.PrintLogo()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)
	select {
	case <-stop:
		(*host).Close()
		os.Exit(0)
	}
}
