package app

import (
	"bufio"
	"fmt"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/services"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/services/random"

	//"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/services"
	"os"

	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/dependencies"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	log "github.com/sirupsen/logrus"
)

// TOD set protocol as random
const Protocol = "/blitzshare/1.0.0"

func readFromStdinToStream(rw *bufio.ReadWriter) {
	for {
		str, _ := rw.ReadString('\n')
		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}
	}
}

func writeStreamFromStdin(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}

		rw.WriteString(fmt.Sprintf("%s\n", sendData))
		rw.Flush()
	}
}

func StartPeer(dep *dependencies.Dependencies) *host.Host {
	h, multiAddr := dep.P2p.StartPeer(dep.Config, Protocol, func(s network.Stream) {
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		go readFromStdinToStream(rw)
		go writeStreamFromStdin(rw)
	})
	otp := random.GenerateRandomWords()
	dep.BlitzshareApi.RegisterAsPeer(multiAddr, otp)

	services.CopyToClipBoard(&otp)
	log.Printf("P2p Address: %s", multiAddr)
	log.Printf("P2p OTP: %s", otp)

	log.Printf("(OTP Copied to Clip Board)")
	return h
}

func ConnectToPeerPass(dep *dependencies.Dependencies, pass *string) *host.Host {
	h := dep.P2p.ConnectToBootsrapNode(dep.Config)
	address := dep.BlitzshareApi.GetPeerAddr(pass)
	log.Printf("[Connecting] OTP: %s", *pass)
	rw := dep.P2p.ConnectToPeer(h, &address.MultiAddr, Protocol)
	log.Printf("[Connected] P2p Address: %s", address.MultiAddr)
	go writeStreamFromStdin(rw)
	go readFromStdinToStream(rw)
	return h
}
