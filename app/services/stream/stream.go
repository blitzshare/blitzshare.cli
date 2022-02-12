package stream

import (
	"bufio"
	"fmt"
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
	readFile, err := os.Open(file)
	defer readFile.Close()
	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		line := fileScanner.Bytes()
		_, err := rw.Write(line)
		_, err = rw.WriteString("\n")
		if err != nil {
			log.Fatalln("failed to write file line to peer stream", err.Error())
		}
	}
	err = rw.Flush()
	// TODO: wait for the stream to finish writing instead of hardcodindg magic numbers
	time.Sleep(time.Second * 5)
	if err == nil {
		log.Println("File sent")
	}
}

func SaveStreamToFile(rw *bufio.ReadWriter, fileName *string) error {
	fHandle, err := os.Create(*fileName)
	if err != nil {
		log.Fatal(err)
	}
	for {
		var line, _, err = rw.ReadLine()
		if err != nil {
			if err.Error() == "stream reset" {
				rw.Flush()
				return nil
			} else {
				log.Panic(err)
			}
		}
		fmt.Fprintln(fHandle, string(line))
	}
	return err
}
