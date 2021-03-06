package app

import (
	"bufio"
	"fmt"
	"strings"

	"os"

	"bootstrap.cli/app/dependencies"
	"bootstrap.cli/app/services/blitzshare"
	"bootstrap.cli/app/services/stream"

	"github.com/libp2p/go-libp2p-core/network"
	log "github.com/sirupsen/logrus"
)

type OTP = string

func StartPeer(dep *dependencies.Dependencies) *OTP {
	mode := "chat"
	otp := dep.Rnd.GenerateRandomWordSequence()
	var token *string
	multiAddr := dep.P2p.StartPeer(dep.Config, otp, func(s network.Stream) {
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		go stream.WriteStreamFromStdin(rw, func() {
			dep.BlitzshareApi.DeregisterAsPeer(otp, token)
		})
		go stream.ReadStreamToStdIo(rw)
	})
	token = dep.BlitzshareApi.RegisterAsPeer(&multiAddr, otp, &mode)
	dep.ClipBoard.CopyToClipBoard(otp)
	return otp
}

func StartPeerFs(dep *dependencies.Dependencies, file string) *OTP {
	mode := "file"
	otp := dep.Rnd.GenerateRandomWordSequence()
	var token *string
	multiAddr := dep.P2p.StartPeer(dep.Config, otp, func(s network.Stream) {
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		dep.BlitzshareApi.DeregisterAsPeer(otp, token)
		fmt.Println("Connected established successfully, sending file")
		stream.SendFileToStream(file, rw)
		os.Exit(0)
	})
	token = dep.BlitzshareApi.RegisterAsPeer(&multiAddr, otp, &mode)
	dep.ClipBoard.CopyToClipBoard(otp)
	return otp
}

func ConnectToPeerOTP(dep *dependencies.Dependencies, otp *string) *blitzshare.P2pPeerRegistryResponse {
	config := dep.BlitzshareApi.GetPeerConfig(otp)
	log.Printf("Connect to peer OTP: %s, mode: %s", *otp, config.Mode)
	log.Printf("Continue? [Y/n]")
	r := bufio.NewReader(os.Stdin)
	s, _ := r.ReadString('\n')
	s = strings.TrimSpace(s)
	if s == "n" {
		ExitProc()
	}
	rw := dep.P2p.ConnectToPeer(dep.Config, &config.MultiAddr, otp)
	log.Printf("[Connected] P2p Address: %s", config.MultiAddr)
	if config.Mode == "chat" {
		go stream.WriteStreamFromStdin(rw, nil)
		go stream.ReadStreamToStdIo(rw)
	} else {
		fileName := fmt.Sprintf("blitzshare-%s.txt", *otp)
		err := stream.SaveStreamToFile(rw, &fileName)
		if err != nil {
			log.Fatalln("error writing to file", fileName, err)
		}
		log.Println("[Success] Peer data written to", fileName)
		ExitProc()
	}
	return config
}

func ExitProc() {
	os.Exit(0)
}
