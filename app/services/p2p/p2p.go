package p2p

import (
	"bufio"
	"context"
	"fmt"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/config"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/services/net"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

type P2p interface {
	StartPeer(conf *config.AppConfig, protocol protocol.ID, handler func(s network.Stream)) string
	ConnectToBootsrapNode(conf *config.AppConfig) *host.Host
	ConnectToPeer(conf *config.AppConfig, address *string, protocol protocol.ID) *bufio.ReadWriter
	Close() error
}

type P2pImp struct {
	host *host.Host
	P2p
}

func NewP2p() P2p {
	return &P2pImp{}
}

func (impl *P2pImp) Close() error {
	return (*impl.host).Close()
}

func (impl *P2pImp) StartPeer(conf *config.AppConfig, protocol protocol.ID, handler func(s network.Stream)) string {
	impl.host = impl.ConnectToBootsrapNode(conf)
	(*impl.host).SetStreamHandler(protocol, handler)
	multiAddr := fmt.Sprintf("/ip4/%s/tcp/%v/p2p/%s \n", conf.LocalP2pPeerIp, net.GetPort(*impl.host), (*impl.host).ID().Pretty())
	return multiAddr
}

func (impl *P2pImp) ConnectToPeer(conf *config.AppConfig, address *string, protocol protocol.ID) *bufio.ReadWriter {
	h := impl.ConnectToBootsrapNode(conf)
	maddr, err := multiaddr.NewMultiaddr(*address)
	if err != nil {
		log.Fatalln(err)
	}
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Fatalln(err)
	}
	(*h).Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
	s, err := (*h).NewStream(context.Background(), info.ID, protocol)
	if err != nil {
		log.Fatalln(err)
	}
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	if err != nil {
		log.Fatalln(err)
	}
	return rw
}

func (*P2pImp) ConnectToBootsrapNode(conf *config.AppConfig) *host.Host {
	log.Printf("[Connecting] P2p network")
	ctx := context.Background()
	host, err := libp2p.New(ctx,
		//libp2p.Security(tls.ID, tls.New),
		libp2p.EnableRelay(),
	)
	targetAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d/p2p/%s", conf.P2pBoostrapNodeIp, conf.P2pBoostrapNodePort, conf.P2pBoostrapNodeId))
	if err != nil {
		log.Fatalln(err)
	}
	targetInfo, err := peer.AddrInfoFromP2pAddr(targetAddr)
	if err != nil {
		log.Fatalln(err)
		log.Fatalln(err)
	}
	err = host.Connect(ctx, *targetInfo)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("[Connected] %s", targetAddr)

	return &host
}
