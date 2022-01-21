package blitzshare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	cfg "bootstrap.cli/app/config"
	log "github.com/sirupsen/logrus"
)

const (
	ChatMode = "chat"
	FileMode = "file"
)

type NodeConfigRespone struct {
	NodeId string `json:"nodeId"`
	Port   int    `json:"port"`
}

type P2pPeerRegistryResponse struct {
	MultiAddr string `form:"multiAddr" binding:"required" json:"multiAddr"`
	Otp       string `form:"otp" binding:"required" json:"otp"`
	Mode      string `form:"otp" binding:"required" json:"mode"`
}

type PeerRegistryAckResponse struct {
	AckId string `binding:"required" json:"ackId"`
	Token string `binding:"required" json:"token"`
}

type BlitzshareApi interface {
	RegisterAsPeer(multiAddr, oneTimePass, mode *string) *string
	GetPeerAddr(oneTimePass *string) *P2pPeerRegistryResponse
	GetBootstrapNode() *NodeConfigRespone
	DeregisterAsPeer(otp, token *string) bool
}

type BlitzshareApiImpl struct {
	BaseUrl string
}

func NewBlitzsahreApi(config *cfg.AppConfig) BlitzshareApi {
	return &BlitzshareApiImpl{BaseUrl: config.BlitzshareApiUrl}
}

func (impl *BlitzshareApiImpl) DeregisterAsPeer(otp, token *string) bool {
	url := fmt.Sprintf("%s/p2p/registry/%s/%s", impl.BaseUrl, *otp, *token)
	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", url, nil)
	response, err := client.Do(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer response.Body.Close()
	log.Debugln(response.StatusCode, url)
	return response.StatusCode == http.StatusAccepted
}

func (impl *BlitzshareApiImpl) RegisterAsPeer(multiAddr *string, otp, mode *string) *string {
	payload, err := json.Marshal(map[string]string{
		"multiAddr": *multiAddr,
		"otp":       *otp,
		"mode":      *mode,
	})
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("%s/p2p/registry", impl.BaseUrl)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	defer resp.Body.Close()
	ack := PeerRegistryAckResponse{}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &ack)
	if err != nil {
		log.Fatalln(err)
	}
	log.Debugln("RegisterAsPeer", ack.Token, ack.AckId, url)
	return &ack.Token
}

func (impl *BlitzshareApiImpl) GetPeerAddr(oneTimePass *string) *P2pPeerRegistryResponse {
	url := fmt.Sprintf("%s/p2p/registry/%s", impl.BaseUrl, *oneTimePass)
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	peerAddress := P2pPeerRegistryResponse{}
	err = json.Unmarshal(body, &peerAddress)
	if err != nil {
		fmt.Println(err)
	}
	return &peerAddress
}

func (impl *BlitzshareApiImpl) GetBootstrapNode() *NodeConfigRespone {
	url := fmt.Sprintf("%s/p2p/bootstrap-node", impl.BaseUrl)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	config := NodeConfigRespone{}
	err = json.Unmarshal(body, &config)
	if err == nil {
		return &config
	}
	return nil
}
