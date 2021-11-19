package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/chat"
	"github.com/libp2p/go-libp2p-core/host"
)

const IP = "127.0.0.1"
const ID = "12D3KooWRZpNwYDJLfErZAVFxXXtomYkkhBrEQCtZx9paWrDk2cd"
const PORT = 63785

func main() {
	dest := flag.String("d", "", "Destination multiaddr string")
	flag.Parse()
	var host host.Host
	conf := &chat.BootstrapP2pConfig{Port: PORT, Ip: IP, NodeId: ID}
	if *dest == "" {
		host = chat.StartPeer(conf)
	} else {
		host = chat.ConnectToPeer(dest, conf)
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)
	select {
	case <-stop:
		host.Close()
		os.Exit(0)
	}
}
