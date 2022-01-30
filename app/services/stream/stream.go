package stream

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func WriteStreamFromStdin(rw *bufio.ReadWriter, onPeerConnectedCb func()) {
	stdReader := bufio.NewReader(os.Stdin)
	onConnectedCalled := false
	for {
		fmt.Print("> ")
		if !onConnectedCalled && onPeerConnectedCb != nil {
			onConnectedCalled = true
			onPeerConnectedCb()
		}
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}
		rw.WriteString(fmt.Sprintf("%s\n", sendData))
		rw.Flush()
	}
}

func ReadStreamToStdIo(rw *bufio.ReadWriter) {
	for {
		str, _ := rw.ReadString('\n')
		if str != "" && str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}
	}
}

func SendFileToStream(file string, rw *bufio.ReadWriter) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln("file cannot be read", file)
	}
	nn, err := rw.Write(content)
	if err != nil {
		log.Fatalln("falied to write file contenct to peer stream", err.Error())
	}
	// err = rw.Flush()
	// TODO: wait for the stream to finish writing instead of hardcodindg magic numbers
	time.Sleep(time.Second * 2)
	if nn == len(content) || err == nil {
		log.Println("File sent")
	}
}
