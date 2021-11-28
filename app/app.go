package app

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/dependencies"
	net "github.com/blitzshare/blitzshare.bootstrap.client.cli/app/services/net"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/services/random"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

// TOD set protocol as random w
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
	words := random.GenerateRandomWords()
	h, err := connectToBootsrapNode(dep)
	if err != nil {
		log.Fatalln(err)
	}
	(*h).SetStreamHandler(Protocol, func(s network.Stream) {
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		go readFromStdinToStream(rw)
		go writeStreamFromStdin(rw)
	})

	multiAddr := fmt.Sprintf("/ip4/%s/tcp/%v/p2p/%s \n", dep.Config.LocalP2pPeerIp, net.GetPort(*h), (*h).ID().Pretty())
	dep.BlitzshareApi.RegisterAsPeer(multiAddr, words)
	log.Printf("P2p Address: %s", multiAddr)
	log.Printf("P2p OTP: %s", words)
	//log.Printf("run: go run ./cmd/*.go -p %s\n", words)
	return h
}

func ConnectToPeerPass(dep *dependencies.Dependencies, pass *string) *host.Host {
	h, err := connectToBootsrapNode(dep)
	if err != nil {
		log.Fatalln(err)
	}
	address := dep.BlitzshareApi.GetPeerAddr(pass)
	log.Printf("[Connecting] OTP: %s", *pass)

	rw := connectToPeer(h, &address.MultiAddr)

	log.Printf("[Connected] P2p Address: %s", address.MultiAddr)

	go writeStreamFromStdin(rw)
	go readFromStdinToStream(rw)
	return h
}

func ConnectToPeerAddress(dep *dependencies.Dependencies, address *string) *host.Host {
	log.Infoln("ConnectToPeerAddress", address)
	h, err := connectToBootsrapNode(dep)
	if err != nil {
		log.Fatalln(err)
	}
	rw := connectToPeer(h, address)

	go writeStreamFromStdin(rw)
	go readFromStdinToStream(rw)

	return h
}

func connectToBootsrapNode(dep *dependencies.Dependencies) (*host.Host, error) {
	log.Printf("[Connecting] P2p network")
	ctx := context.Background()
	host, err := libp2p.New(ctx,
		//libp2p.Security(tls.ID, tls.New),
		libp2p.EnableRelay(),
	)
	targetAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d/p2p/%s", dep.Config.P2pBoostrapNodeIp, dep.Config.P2pBoostrapNodePort, dep.Config.P2pBoostrapNodeId))
	if err != nil {
		log.Panicln(err)
	}
	targetInfo, err := peer.AddrInfoFromP2pAddr(targetAddr)
	if err != nil {
		log.Panicln(err)
	}
	err = host.Connect(ctx, *targetInfo)
	if err != nil {
		log.Panicln(err)
	}
	log.Printf("[Connected] %s", targetAddr)

	return &host, err
}

func connectToPeer(h *host.Host, destination *string) *bufio.ReadWriter {
	maddr, err := multiaddr.NewMultiaddr(*destination)
	if err != nil {
		log.Fatalln(err)
	}
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Fatalln(err)
	}
	(*h).Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
	s, err := (*h).NewStream(context.Background(), info.ID, Protocol)
	if err != nil {
		log.Fatalln(err)
	}
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	if err != nil {
		log.Fatalln(err)
	}
	return rw
}
