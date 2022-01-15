package app

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"os"

	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/dependencies"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/services/blitzshare"

	"github.com/libp2p/go-libp2p-core/network"
	log "github.com/sirupsen/logrus"
)

type OTP = string

func readFromStdinToStream(rw *bufio.ReadWriter) {
	for {
		str, _ := rw.ReadString('\n')
		if str != "" && str != "\n" {
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

func StartPeer(dep *dependencies.Dependencies) *OTP {
	mode := "chat"
	otp := dep.Rnd.GenerateRandomWordSequence()
	multiAddr := dep.P2p.StartPeer(dep.Config, otp, func(s network.Stream) {
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		startChat(rw)
	})
	dep.BlitzshareApi.RegisterAsPeer(&multiAddr, otp, &mode)
	dep.ClipBoard.CopyToClipBoard(otp)
	return otp
}
func SendFileToStream(file string, rw *bufio.ReadWriter) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln("file %s cannot be read", file)
	}
	_, err = rw.Write(content)
	if err != nil {
		log.Fatalln("falied to write file contenct to peer stream", err.Error())
	}
	err = rw.Flush()
	// TODO: wait for the stream to finish writing instead of hardcodindg magic numbers
	time.Sleep(time.Second * 5)
	if err == nil {
		log.Println("File sent")
	} else {
		log.Fatalln("falied to write file contenct to peer stream", err.Error())
	}
}

func StartPeerFs(dep *dependencies.Dependencies, file string) *OTP {
	mode := "file"
	otp := dep.Rnd.GenerateRandomWordSequence()
	var token *string
	multiAddr := dep.P2p.StartPeer(dep.Config, otp, func(s network.Stream) {
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		dep.BlitzshareApi.DeregisterAsPeer(otp, token)
		fmt.Println("Connected established successfully, sending file")
		SendFileToStream(file, rw)

		os.Exit(0)
	})
	token = dep.BlitzshareApi.RegisterAsPeer(&multiAddr, otp, &mode)
	dep.ClipBoard.CopyToClipBoard(otp)
	return otp
}

func startChat(rw *bufio.ReadWriter) {
	go writeStreamFromStdin(rw)
	go readFromStdinToStream(rw)
}

func SaveStreamToFile(rw *bufio.ReadWriter, otp *string) {
	bytes, err := ioutil.ReadAll(rw)
	if err == nil {
		log.Fatalln("faield to receive file from peer stream")
	} else {
		fileName := fmt.Sprintf("blitzshare-%s.txt", *otp)
		if err := os.WriteFile(fileName, bytes, 0666); err != nil {
			log.Fatal(err)
		}
		log.Printf("file saved as %s", fileName)
	}
}
func ConnectToPeerOTP(dep *dependencies.Dependencies, otp *string) *blitzshare.P2pPeerRegistryResponse {
	config := dep.BlitzshareApi.GetPeerAddr(otp)
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
		startChat(rw)
	} else {
		SaveStreamToFile(rw, otp)
		ExitProc()
	}
	return config
}

func ExitProc() {
	os.Exit(0)
}
