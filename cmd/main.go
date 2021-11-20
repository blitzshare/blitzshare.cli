package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/chat"
	"github.com/libp2p/go-libp2p-core/host"
)

const IP = "0.0.0.0"
const ID = "12D3KooWAZK68L58pAoqQfCtRFUPU2YE4A6Ang8D44QEHqnJxWz5"
const PORT = 63785

func main() {
	dest := flag.String("d", "", "Destination multiaddr string")
	pass := flag.String("p", "", "One time pass id of connection peer")
	flag.Parse()
	var host host.Host
	conf := &chat.BootstrapP2pConfig{Port: PORT, Ip: IP, NodeId: ID}
	if *dest != "" {
		host = chat.ConnectToPeerAddress(conf, dest)
	} else if *pass != "" {
		host = chat.ConnectToPeerPass(conf, pass)
	} else {
		host = chat.StartPeer(conf)
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)
	select {
	case <-stop:
		host.Close()
		os.Exit(0)
	}
}
