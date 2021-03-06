package p2p

import (
	"bufio"
	"context"
	"fmt"

	"bootstrap.cli/app/config"
	"bootstrap.cli/app/services/str"
	"github.com/libp2p/go-libp2p"
	tls "github.com/libp2p/go-libp2p-tls"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

type P2p interface {
	StartPeer(conf *config.AppConfig, otp *string, handler func(s network.Stream)) string
	ConnectToBootsrapNode(conf *config.AppConfig) *host.Host
	ConnectToPeer(conf *config.AppConfig, address *string, otp *string) *bufio.ReadWriter
	Close() error
}

type P2pImp struct {
	host *host.Host
}

func NewP2p() P2p {
	return &P2pImp{}
}

func (impl *P2pImp) Close() error {
	return (*impl.host).Close()
}

func getProtocol(otp *string) protocol.ID {
	proto := fmt.Sprintf("/blitzshare/1.0.0/%s", *otp)
	return protocol.ID(proto)
}

func (impl *P2pImp) StartPeer(conf *config.AppConfig, otp *string, handler func(s network.Stream)) string {
	impl.host = impl.ConnectToBootsrapNode(conf)
	(*impl.host).SetStreamHandler(getProtocol(otp), handler)
	multiAddr := fmt.Sprintf("/ip4/%s/tcp/%v/p2p/%s", conf.LocalP2pPeerIp, GetPort(*impl.host), (*impl.host).ID().Pretty())
	return multiAddr
}

func (impl *P2pImp) ConnectToPeer(conf *config.AppConfig, address, otp *string) *bufio.ReadWriter {
	h := impl.ConnectToBootsrapNode(conf)
	maddr, err := multiaddr.NewMultiaddr(str.SanatizeStr(*address))
	if err != nil {
		log.Fatalln(err)
	}
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Fatalln(err)
	}
	(*h).Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
	s, err := (*h).NewStream(context.Background(), info.ID, getProtocol(otp))
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
	h, err := libp2p.New(ctx,
		libp2p.Security(tls.ID, tls.New),
		libp2p.EnableRelay(),
	)
	targetAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d/p2p/%s", str.SanatizeStr(conf.P2pBoostrapNodeIp), conf.P2pBoostrapNodePort, str.SanatizeStr(conf.P2pBoostrapNodeId)))
	if err != nil {
		log.Fatalln(err)
	}
	targetInfo, err := peer.AddrInfoFromP2pAddr(targetAddr)
	if err != nil {
		log.Fatalln(err)
		log.Fatalln(err)
	}
	err = h.Connect(ctx, *targetInfo)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("[Connected] %s", targetAddr)

	return &h
}
