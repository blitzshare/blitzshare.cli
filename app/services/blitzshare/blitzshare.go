package blitzshare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	cfg "github.com/blitzshare/blitzshare.bootstrap.client.cli/app/config"
	log "github.com/sirupsen/logrus"
)

type PeerAddress struct {
	MultiAddr string `json:"multiAddr"`
}

type BlitzshareApi interface {
	RegisterAsPeer(multiAddr string, oneTimePass *string) bool
	GetPeerAddr(oneTimePass *string) *PeerAddress
}

type BlitzshareApiImpl struct {
	BaseUrl string
}

func NewBlitzsahreApi(config *cfg.AppConfig) BlitzshareApi {
	return &BlitzshareApiImpl{BaseUrl: config.BlitzshareApiUrl}
}

func (impl *BlitzshareApiImpl) RegisterAsPeer(multiAddr string, oneTimePass *string) bool {
	body, err := json.Marshal(map[string]string{
		"multiAddr": multiAddr,
		"otp":       *oneTimePass,
	})
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("%s/p2p/registry", impl.BaseUrl)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}
	return resp.Status == "202 Accepted"
}

func (impl *BlitzshareApiImpl) GetPeerAddr(oneTimePass *string) *PeerAddress {
	url := fmt.Sprintf("%s/p2p/registry/%s", impl.BaseUrl, *oneTimePass)
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	peerAddress := PeerAddress{}
	err = json.Unmarshal(body, &peerAddress)
	if err != nil {
		fmt.Println(err)
	}
	return &peerAddress
}
