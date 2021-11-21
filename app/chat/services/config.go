package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const BASE_CONFIG_API_URL = "http://0.0.0.0"

type PeerAddress struct {
	MultiAddr string `json:"multiAddr"`
}

func RegisterAsPeer(multiAddr string, oneTimePass string) bool {
	body, err := json.Marshal(map[string]string{
		"multiAddr":   multiAddr,
		"oneTimePass": oneTimePass,
	})
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("%s/p2p/registry", BASE_CONFIG_API_URL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}
	return resp.Status == "200 OK"
}

func GetPeerAddr(oneTimePass *string) *PeerAddress {
	url := fmt.Sprintf("%s/p2p/registry/%s", BASE_CONFIG_API_URL, *oneTimePass)
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	peerAddress := PeerAddress{}
	err = json.Unmarshal(body, &peerAddress)
	if err != nil {
		fmt.Println(err)
	}
	return &peerAddress
}
