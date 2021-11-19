package pubsub

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-libp2p-core/routing"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	mplex "github.com/libp2p/go-libp2p-mplex"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	tls "github.com/libp2p/go-libp2p-tls"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/libp2p/go-tcp-transport"
	ws "github.com/libp2p/go-ws-transport"
	"github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

type mdnsNotifee struct {
	h   host.Host
	ctx context.Context
	mdns.Notifee
}

func (m *mdnsNotifee) HandlePeerFound(pi peer.AddrInfo) {
	m.h.Connect(m.ctx, pi)
}

const pubsubTopic = "/chat/1.0.0" // "test-chat-topic"
var MyId string

func StartPubSub() {
	argsWithoutProg := os.Args[1:]
	sender := argsWithoutProg[0]
	// ID := argsWithoutProg[1]
	IP := "10.101.18.26" // "3.10.185.121"
	ID := "12D3KooWMr8ABosc5unTHDsCN1QBYAbvmWWppH53tvtMJrqM3kFC"

	log.Infoln("ID", ID)
	log.Infoln("IP", IP)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	transports := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(ws.New),
	)
	fmt.Print(transports)

	muxers := libp2p.ChainOptions(
		// libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
		libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport),
	)

	security := libp2p.Security(tls.ID, tls.New)

	listenAddrs := libp2p.ListenAddrStrings(
		"/ip4/0.0.0.0/tcp/0",
		//"/ip4/0.0.0.0/tcp/0/ws",
		//"/ip4/127.0.0.1/tcp/0",
	)

	var dht *kaddht.IpfsDHT
	newDHT := func(h host.Host) (routing.PeerRouting, error) {
		var err error
		dht, err = kaddht.New(ctx, h)
		return dht, err
	}
	routing := libp2p.Routing(newDHT)

	host, err := libp2p.New(ctx,
		transports,
		listenAddrs,
		muxers,
		security,
		routing,
		libp2p.EnableRelay(),
		libp2p.NATPortMap(),
	)
	MyId = host.ID().String()
	log.Infoln("[I AM]", MyId)
	if sender != "sender" {
		registerStreams(host)
	}
	if err != nil {
		panic(err)
	}

	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(err)
	}
	topic, err := ps.Join(pubsubTopic)
	if err != nil {
		panic(err)
	}
	defer topic.Close()
	sub, err := topic.Subscribe()
	if err != nil {
		panic(err)
	}
	go pubsubHandler(ctx, sub)

	for _, addr := range host.Addrs() {
		fmt.Println("Listening on", addr)
	}

	if sender == "sender" {
		go iteratePeersInChannel(ctx, host, topic)
		go chatInputLoop(ctx, host, topic)
	} else {
		go chatInputLoop(ctx, host, topic)
	}

	connection := "/ip4/" + IP + "/tcp/63785/p2p/" + ID
	targetAddr, err := multiaddr.NewMultiaddr(connection)
	if err != nil {
		panic(err)
	}
	targetInfo, err := peer.AddrInfoFromP2pAddr(targetAddr)
	if err != nil {
		panic(err)
	}

	err = host.Connect(ctx, *targetInfo)
	if err != nil {
		panic(err)
	}
	log.Infoln("Connected to", targetAddr)

	sa := mdns.NewMdnsService(host, "findme") // a, time.Second/4,
	err = dht.Bootstrap(ctx)
	if err != nil {
		panic(err)
	}
	notifee := &mdnsNotifee{h: host, ctx: ctx}
	sa.RegisterNotifee(notifee)
	// fmt.Println(sa)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)

	// if sender != "sender" {
	// 	registerStreams(host)
	// }
	select {
	case <-stop:
		host.Close()
		os.Exit(0)
	}
}

func registerStreams(host host.Host) {
	log.Infoln("Set Stream handler", protocol.ID(ProtocolID))

	host.SetStreamHandler(protocol.ID(ProtocolID), func(stream network.Stream) {
		log.Info("@@ SetStreamHandler Got a new stream!")
	})

	host.SetStreamHandlerMatch(protocol.ID(ProtocolID), func(s string) bool {
		log.Info("@@ SetStreamHandlerMatch ", s)
		return true
	}, func(stream network.Stream) {
		log.Info("@@ SetStreamHandlerMatch Got a new stream! ", stream.Protocol())
		// rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
		// for {
		// 	str, err := rw.ReadString('\n')
		// 	if err != nil {
		// 		fmt.Println("Error reading from buffer", err)
		// 	}
		// 	log.Info("@@ Got message ", str)

		// }
	})
	host.Network().SetConnHandler(func(con network.Conn) {
		// log.Infoln("@@ Peer Joined ", con.RemotePeer(), con.LocalMultiaddr().String())
	})
}

const ProtocolID = "p2p/protocol"

var peersInTopic []string

func sendMessageWithDirectStream(ctx context.Context, h host.Host, p peer.ID) {
	s, err := h.NewStream(ctx, p, protocol.ID(ProtocolID))
	log.Infoln("remote-peer", p.Pretty())
	log.Infoln("remote-peer:protocol", s.Protocol())
	log.Infoln("remote-peer:stat", s.Stat())
	s.SetProtocol(protocol.ID(ProtocolID))
	if err != nil {
		log.Errorln("Error dial peer", err)
	}

	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	msgId := make([]byte, 10)
	_, _ = rand.Read(msgId)
	now := time.Now().Unix()
	req := &Request{
		Type: Request_SEND_MESSAGE.Enum(),
		SendMessage: &SendMessage{
			Id:      msgId,
			Data:    []byte("hi there"),
			Created: &now,
		},
	}
	msgBytes, err := proto.Marshal(req)
	if err != nil {
		log.Errorln("Marshal req", err)
	}
	n, err := rw.Write(msgBytes)
	log.Infoln("n", n, err)
	log.Infoln("direct message sent to peer", p.String())
}

func onNewPeerFound(ctx context.Context, h host.Host, p peer.ID, addr []multiaddr.Multiaddr) {
	// con := h.Network().ConnsToPeer(p)
	// for _, c := range con {
	// 	streams := c.GetStreams()
	// 	for _, s := range streams {
	// 		log.Infoln("remote-peer", c.RemotePeer().Pretty())
	// 		log.Infoln("remote-peer:protocol", s.Protocol())

	// 		msgId := make([]byte, 10)
	// 		_, _ = rand.Read(msgId)
	// 		now := time.Now().Unix()
	// 		req := &Request{
	// 			Type: Request_SEND_MESSAGE.Enum(),
	// 			SendMessage: &SendMessage{
	// 				Id:      msgId,
	// 				Data:    []byte("hi there"),
	// 				Created: &now,
	// 			},
	// 		}
	// 		msgBytes, err := proto.Marshal(req)
	// 		if err != nil {
	// 			log.Errorln("Marshal req", err)
	// 		}
	// 		n, err := s.Write(msgBytes)
	// 		log.Infoln("direct message sent to peer", p.String())
	// 		log.Infoln("n", n, err)
	// 	}
	// }
	//for _, ad := range addr {
	//
	//	m, _ := multiaddr.NewMultiaddr(ad.String() + "/" + p.String())
	//	peerinfo, _ := peer.AddrInfoFromP2pAddr(m)
	//	p2pm, _ := multiaddr.NewMultiaddr("p2p/" + p.String())
	//	p2pPeerinfo, _ := peer.AddrInfoFromP2pAddr(p2pm)
	//	fmt.Println("p2p peerinfo", p2pPeerinfo)
	//	fmt.Println("info", peerinfo)
	//	//if peerinfo != nil {
	//	//	fmt.Println("info", peerinfo)
	//	//	if err := h.Connect(ctx, *peerinfo); err != nil {
	//	//		fmt.Println(err)
	//	//	} else {
	//	//		fmt.Println("Connection established with bootstrap node: ", *peerinfo)
	//	//	}
	//	//} else {
	//	//	fmt.Println("info is nil")
	//	//}
	//}
	// /ip4/1.2.3.4/tcp/4321/p2p/QmcEPrat8ShnCph8WjkREzt5CPXF2RwhYxYBALDcLC1iV6

	//peerinfo, _ := peer.AddrInfoFromP2pAddr(ad)
	//addr, _ := multiaddr.NewMultiaddr(ad.String())

	//peersInTopic = append(peersInTopic, p.String())
	//// TODO: understand why we get - ERRO[0005] Error dial peer protocol not supported
	//log.Infoln("New Peer", p.String())
}

func iteratePeersInChannel(ctx context.Context, h host.Host, topic *pubsub.Topic) {
	for {
		store := h.Peerstore()
		peers := store.Peers()
		if len(peers) > 0 {
			for _, p := range peers { //topic.ListPeers()
				// log.Infoln("peers", peers)
				if p.String() != MyId {
					// sendMessageWithDirectStream(ctx, h, p)
					con, err := h.Network().DialPeer(ctx, p)
					if err == nil {
						log.Infoln(con)
						streams := con.GetStreams()
						for _, s := range streams {
							// log.Infoln("remote-peer", con.RemotePeer().Pretty())
							log.Infoln("remote-peer:protocol", s.Protocol())
							rs := bufio.NewWriter(bufio.NewWriter(s))
							rs.WriteString("hello to you from " + MyId)
							err = rs.Flush()
							if err != nil {
								log.Errorln("Flush err", err)
							}
							log.Infoln("direct message sent to peer", p.String())
						}

					} else {
						log.Errorln("faield dial peer", err)

					}
				}

			}

		} else {
			log.Infoln("no peers found :(")
		}
		time.Sleep(5 * time.Second)
	}
}

func sendMessage(ctx context.Context, topic *pubsub.Topic, msg string) {
	msgId := make([]byte, 10)
	_, err := rand.Read(msgId)
	defer func() {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()
	if err != nil {
		return
	}
	now := time.Now().Unix()
	req := &Request{
		Type: Request_SEND_MESSAGE.Enum(),
		SendMessage: &SendMessage{
			Id:      msgId,
			Data:    []byte(msg),
			Created: &now,
		},
	}
	msgBytes, err := proto.Marshal(req)
	if err != nil {
		return
	}
	err = topic.Publish(ctx, msgBytes)
}

var handles = map[string]string{}

func updatePeer(ctx context.Context, topic *pubsub.Topic, id peer.ID, handle string) {
	oldHandle, ok := handles[id.String()]
	if !ok {
		oldHandle = id.ShortString()
	}
	handles[id.String()] = handle

	req := &Request{
		Type: Request_UPDATE_PEER.Enum(),
		UpdatePeer: &UpdatePeer{
			UserHandle: []byte(handle),
		},
	}
	reqBytes, err := proto.Marshal(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	err = topic.Publish(ctx, reqBytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Printf("%s -> %s\n", oldHandle, handle)
}

func chatInputLoop(ctx context.Context, h host.Host, topic *pubsub.Topic) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		sendMessage(ctx, topic, "test message")
	}
}

func pubsubHandler(ctx context.Context, sub *pubsub.Subscription) {
	defer sub.Cancel()
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		log.Infoln("@ msg @", string(msg.Data))
		req := &Request{}
		err = proto.Unmarshal(msg.Data, req)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		switch *req.Type {
		case Request_SEND_MESSAGE:
			pubsubMessageHandler(msg.GetFrom(), req.SendMessage)
		case Request_UPDATE_PEER:
			pubsubUpdateHandler(msg.GetFrom(), req.UpdatePeer)
		}
	}
}
func pubsubMessageHandler(id peer.ID, msg *SendMessage) {
	handle, ok := handles[id.String()]
	if !ok {
		handle = id.ShortString()
	}
	fmt.Printf("%s: %s\n", handle, msg.Data)
}

func pubsubUpdateHandler(id peer.ID, msg *UpdatePeer) {
	oldHandle, ok := handles[id.String()]
	if !ok {
		oldHandle = id.ShortString()
	}
	handles[id.String()] = string(msg.UserHandle)
	fmt.Printf("%s -> %s\n", oldHandle, msg.UserHandle)
}
