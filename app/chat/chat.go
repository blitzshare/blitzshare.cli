package chat

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

const Protocol = "/blitzshare/chat/1.0.0"

func handleStream(s network.Stream) {
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go readData(rw)
	go writeData(rw)
}

func readData(rw *bufio.ReadWriter) {
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

func writeData(rw *bufio.ReadWriter) {
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

const IP = "10.101.18.26"
const ID = "12D3KooWMr8ABosc5unTHDsCN1QBYAbvmWWppH53tvtMJrqM3kFC"
const PORT = 63785

func StartPeer() host.Host {
	words := GenerateRandomWords()
	log.Infoln(words)
	h, err := connectToBootsrapNode(IP, ID, PORT)
	if err != nil {
		log.Fatalln(err)
	}
	h.SetStreamHandler(Protocol, handleStream)
	log.Printf("Connect Peer: go run ./cmd/*.go -d /ip4/127.0.0.1/tcp/%v/p2p/%s \n", GetPort(h), h.ID().Pretty())
	return h
}

func ConnectToPeer(dest *string) host.Host {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	h, err := connectToBootsrapNode(IP, ID, PORT)
	if err != nil {
		log.Fatalln(err)
	}
	rw, err := connectToPeer(ctx, h, *dest)
	if err != nil {
		log.Fatalln(err)
	}
	go writeData(rw)
	go readData(rw)

	return h
}

func connectToBootsrapNode(ip string, id string, port int) (host.Host, error) {
	ctx, _ := context.WithCancel(context.Background())
	host, err := libp2p.New(ctx,
		// TODO libp2p.Security(tls.ID, tls.New),
		libp2p.EnableRelay(),
	)
	targetAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d/p2p/%s", ip, port, id))
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
	log.Infoln("Bootsrap Node Address: ", targetAddr)
	return host, err
}

func connectToPeer(ctx context.Context, h host.Host, destination string) (*bufio.ReadWriter, error) {
	maddr, err := multiaddr.NewMultiaddr(destination)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// Add the destination's peer multiaddress in the peerstore.
	h.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
	s, err := h.NewStream(context.Background(), info.ID, Protocol)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	log.Println("Connected to ", destination)
	return rw, nil
}
