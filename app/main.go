package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	mplex "github.com/libp2p/go-libp2p-mplex"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	tls "github.com/libp2p/go-libp2p-tls"
	yamux "github.com/libp2p/go-libp2p-yamux"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/libp2p/go-tcp-transport"
	ws "github.com/libp2p/go-ws-transport"
	"github.com/multiformats/go-multiaddr"
	math "math"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type mdnsNotifee struct {
	h   host.Host
	ctx context.Context
}

func (m *mdnsNotifee) HandlePeerFound(pi peer.AddrInfo) {
	m.h.Connect(m.ctx, pi)
}

const ID = "12D3KooWJ23BYF9jBM2qSoo49fizpevGwdf7anVfHn2L48FodGSA"
const IP = "127.0.0.1" // "172.17.0.2" //
const pubsubTopic = "/libp2p/example/chat/1.0.0"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	transports := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(ws.New),
	)
	fmt.Print(transports)

	muxers := libp2p.ChainOptions(
		libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
		libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport),
	)

	security := libp2p.Security(tls.ID, tls.New)

	listenAddrs := libp2p.ListenAddrStrings(
		"/ip4/0.0.0.0/tcp/0",
		"/ip4/0.0.0.0/tcp/0/ws",
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
	)
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
	connection := "/ip4/" + IP + "/tcp/63785/p2p/" + ID
	targetAddr, err := multiaddr.NewMultiaddr(connection)
	fmt.Println("Connected to", connection)
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

	fmt.Println("Connected to", targetInfo.ID)

	notifee := &mdnsNotifee{h: host, ctx: ctx}
	sa := mdns.NewMdnsService(host, "findme") // a, time.Second/4,
	err = dht.Bootstrap(ctx)
	if err != nil {
		panic(err)
	}
	sa.RegisterNotifee(notifee)
	fmt.Println(sa)

	donec := make(chan struct{}, 1)
	go chatInputLoop(ctx, host, topic, donec)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)

	select {
	case <-stop:
		host.Close()
		os.Exit(0)
	case <-donec:
		host.Close()
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

func chatInputLoop(ctx context.Context, h host.Host, topic *pubsub.Topic, donec chan struct{}) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		if strings.HasPrefix(msg, "/name ") {
			newHandle := strings.TrimPrefix(msg, "/name ")
			newHandle = strings.TrimSpace(newHandle)
			updatePeer(ctx, topic, h.ID(), newHandle)
		} else {
			sendMessage(ctx, topic, msg)
		}
	}
	donec <- struct{}{}
}

func pubsubHandler(ctx context.Context, sub *pubsub.Subscription) {
	defer sub.Cancel()
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

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

//---proto

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type Request_Type int32

const (
	Request_SEND_MESSAGE Request_Type = 0
	Request_UPDATE_PEER  Request_Type = 1
)

var Request_Type_name = map[int32]string{
	0: "SEND_MESSAGE",
	1: "UPDATE_PEER",
}

var Request_Type_value = map[string]int32{
	"SEND_MESSAGE": 0,
	"UPDATE_PEER":  1,
}

func (x Request_Type) Enum() *Request_Type {
	p := new(Request_Type)
	*p = x
	return p
}

func (x Request_Type) String() string {
	return proto.EnumName(Request_Type_name, int32(x))
}

func (x *Request_Type) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Request_Type_value, data, "Request_Type")
	if err != nil {
		return err
	}
	*x = Request_Type(value)
	return nil
}

func (Request_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8c585a45e2093e54, []int{0, 0}
}

type Request struct {
	Type                 *Request_Type `protobuf:"varint,1,req,name=type,enum=main.Request_Type" json:"type,omitempty"`
	SendMessage          *SendMessage  `protobuf:"bytes,2,opt,name=sendMessage" json:"sendMessage,omitempty"`
	UpdatePeer           *UpdatePeer   `protobuf:"bytes,3,opt,name=updatePeer" json:"updatePeer,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_8c585a45e2093e54, []int{0}
}
func (m *Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request.Unmarshal(m, b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request.Marshal(b, m, deterministic)
}
func (m *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(m, src)
}
func (m *Request) XXX_Size() int {
	return xxx_messageInfo_Request.Size(m)
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

func (m *Request) GetType() Request_Type {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return Request_SEND_MESSAGE
}

func (m *Request) GetSendMessage() *SendMessage {
	if m != nil {
		return m.SendMessage
	}
	return nil
}

func (m *Request) GetUpdatePeer() *UpdatePeer {
	if m != nil {
		return m.UpdatePeer
	}
	return nil
}

type SendMessage struct {
	Data                 []byte   `protobuf:"bytes,1,req,name=data" json:"data,omitempty"`
	Created              *int64   `protobuf:"varint,2,req,name=created" json:"created,omitempty"`
	Id                   []byte   `protobuf:"bytes,3,req,name=id" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendMessage) Reset()         { *m = SendMessage{} }
func (m *SendMessage) String() string { return proto.CompactTextString(m) }
func (*SendMessage) ProtoMessage()    {}
func (*SendMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_8c585a45e2093e54, []int{1}
}
func (m *SendMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendMessage.Unmarshal(m, b)
}
func (m *SendMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendMessage.Marshal(b, m, deterministic)
}
func (m *SendMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendMessage.Merge(m, src)
}
func (m *SendMessage) XXX_Size() int {
	return xxx_messageInfo_SendMessage.Size(m)
}
func (m *SendMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_SendMessage.DiscardUnknown(m)
}

var xxx_messageInfo_SendMessage proto.InternalMessageInfo

func (m *SendMessage) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *SendMessage) GetCreated() int64 {
	if m != nil && m.Created != nil {
		return *m.Created
	}
	return 0
}

func (m *SendMessage) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

type UpdatePeer struct {
	UserHandle           []byte   `protobuf:"bytes,1,opt,name=userHandle" json:"userHandle,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdatePeer) Reset()         { *m = UpdatePeer{} }
func (m *UpdatePeer) String() string { return proto.CompactTextString(m) }
func (*UpdatePeer) ProtoMessage()    {}
func (*UpdatePeer) Descriptor() ([]byte, []int) {
	return fileDescriptor_8c585a45e2093e54, []int{2}
}
func (m *UpdatePeer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdatePeer.Unmarshal(m, b)
}
func (m *UpdatePeer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdatePeer.Marshal(b, m, deterministic)
}
func (m *UpdatePeer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdatePeer.Merge(m, src)
}
func (m *UpdatePeer) XXX_Size() int {
	return xxx_messageInfo_UpdatePeer.Size(m)
}
func (m *UpdatePeer) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdatePeer.DiscardUnknown(m)
}

var xxx_messageInfo_UpdatePeer proto.InternalMessageInfo

func (m *UpdatePeer) GetUserHandle() []byte {
	if m != nil {
		return m.UserHandle
	}
	return nil
}

func init() {
	proto.RegisterEnum("main.Request_Type", Request_Type_name, Request_Type_value)
	proto.RegisterType((*Request)(nil), "main.Request")
	proto.RegisterType((*SendMessage)(nil), "main.SendMessage")
	proto.RegisterType((*UpdatePeer)(nil), "main.UpdatePeer")
}

func init() { proto.RegisterFile("chat.proto", fileDescriptor_8c585a45e2093e54) }

var fileDescriptor_8c585a45e2093e54 = []byte{
	// 251 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0x51, 0x4b, 0xc3, 0x30,
	0x14, 0x85, 0x4d, 0x5a, 0x18, 0xdc, 0x96, 0x59, 0xef, 0x53, 0x9e, 0xa4, 0xf4, 0x41, 0x2a, 0x48,
	0x91, 0xf9, 0x0b, 0x06, 0x0b, 0x0a, 0x32, 0x29, 0xe9, 0xf6, 0x3c, 0xc2, 0x72, 0xd1, 0x82, 0x76,
	0xb5, 0x49, 0x1f, 0xf6, 0xfb, 0xfc, 0x63, 0xd2, 0x4c, 0xb7, 0xbc, 0x25, 0xe7, 0x7c, 0x87, 0x73,
	0xef, 0x05, 0xd8, 0x7f, 0x68, 0x57, 0xf5, 0xc3, 0xc1, 0x1d, 0x30, 0xfe, 0xd2, 0x6d, 0x57, 0xfc,
	0x30, 0x98, 0x29, 0xfa, 0x1e, 0xc9, 0x3a, 0xbc, 0x83, 0xd8, 0x1d, 0x7b, 0x12, 0x2c, 0xe7, 0xe5,
	0x7c, 0x81, 0xd5, 0x04, 0x54, 0x7f, 0x66, 0xb5, 0x39, 0xf6, 0xa4, 0xbc, 0x8f, 0x4f, 0x90, 0x58,
	0xea, 0xcc, 0x9a, 0xac, 0xd5, 0xef, 0x24, 0x78, 0xce, 0xca, 0x64, 0x71, 0x73, 0xc2, 0x9b, 0x8b,
	0xa1, 0x42, 0x0a, 0x1f, 0x01, 0xc6, 0xde, 0x68, 0x47, 0x35, 0xd1, 0x20, 0x22, 0x9f, 0xc9, 0x4e,
	0x99, 0xed, 0x59, 0x57, 0x01, 0x53, 0xdc, 0x43, 0x3c, 0x95, 0x62, 0x06, 0x69, 0x23, 0xdf, 0x56,
	0xbb, 0xb5, 0x6c, 0x9a, 0xe5, 0xb3, 0xcc, 0xae, 0xf0, 0x1a, 0x92, 0x6d, 0xbd, 0x5a, 0x6e, 0xe4,
	0xae, 0x96, 0x52, 0x65, 0xac, 0x78, 0x85, 0x24, 0x28, 0x46, 0x84, 0xd8, 0x68, 0xa7, 0xfd, 0x22,
	0xa9, 0xf2, 0x6f, 0x14, 0x30, 0xdb, 0x0f, 0xa4, 0x1d, 0x19, 0xc1, 0x73, 0x5e, 0x46, 0xea, 0xff,
	0x8b, 0x73, 0xe0, 0xad, 0x11, 0x91, 0x67, 0x79, 0x6b, 0x8a, 0x07, 0x80, 0xcb, 0x44, 0x78, 0x0b,
	0x30, 0x5a, 0x1a, 0x5e, 0x74, 0x67, 0x3e, 0xa7, 0xd3, 0xb0, 0x32, 0x55, 0x81, 0xf2, 0x1b, 0x00,
	0x00, 0xff, 0xff, 0x5c, 0xd9, 0x58, 0xd2, 0x53, 0x01, 0x00, 0x00,
}
