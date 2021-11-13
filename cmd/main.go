package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/chat"
	"github.com/libp2p/go-libp2p-core/host"
)

func main() {
	dest := flag.String("d", "", "Destination multiaddr string")
	flag.Parse()
	var host host.Host
	log.Println("dest", string(*dest))
	if *dest == "" {
		host = chat.StartPeer()
	} else {
		host = chat.ConnectToPeer(dest)
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)
	select {
	case <-stop:
		host.Close()
		os.Exit(0)
	}
}
