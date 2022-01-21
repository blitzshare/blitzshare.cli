package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"bootstrap.cli/app/services"
	svc "bootstrap.cli/app/services"
	"bootstrap.cli/app/services/str"

	"bootstrap.cli/app"
	"bootstrap.cli/app/config"
	"bootstrap.cli/app/dependencies"
	log "github.com/sirupsen/logrus"
)

func initLog() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	initLog()
	svc.PrintLogo()
	start := flag.Bool("start", false, "Start p2p init peer session")
	var file string
	flag.StringVar(&file, "file", "", "Start p2p init peer file share session")
	connect := flag.Bool("connect", false, "Start p2p receiver peer session")
	flag.Parse()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config %v\n", err)
	}
	deps, err := dependencies.NewDependencies(cfg)
	if err != nil {
		log.Fatalf("failed to load dependencies %v\n", err)
	}
	node := deps.BlitzshareApi.GetBootstrapNode()
	if node == nil {
		log.Fatalln("Could not get initial node configuration from remote server")
	}
	cfg.P2pBoostrapNodePort = node.Port
	cfg.P2pBoostrapNodeId = node.NodeId
	if *start {
		var otp *string
		if file == "" {
			otp = app.StartPeer(deps)
		} else {
			otp = app.StartPeerFs(deps, file)
		}
		log.Printf("OTP: %s (copied to clipboard)", *otp)
	} else if *connect {
		log.Println("Enter OTP:")
		line := services.ReadStdInLine()
		otp := str.SanatizeStr(*line)
		app.ConnectToPeerOTP(deps, &otp)
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)
	select {
	case <-stop:
		deps.P2p.Close()
		os.Exit(0)
	}
}
