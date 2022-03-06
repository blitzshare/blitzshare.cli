package blitzshare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

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
	GetPeerConfig(oneTimePass *string) *P2pPeerRegistryResponse
	GetBootstrapNode() *NodeConfigRespone
	DeregisterAsPeer(otp, token *string) bool
}

type BlitzshareApiImpl struct {
	baseUrl string
	apiKey  string
}

func NewBlitzsahreApi(config *cfg.AppConfig) BlitzshareApi {
	return &BlitzshareApiImpl{
		baseUrl: config.BlitzshareApiUrl,
		apiKey:  config.BlitzshareApiKey,
	}
}

var client = &http.Client{
	Timeout: time.Second * 10,
}

func (impl *BlitzshareApiImpl) DeregisterAsPeer(otp, token *string) bool {
	url := fmt.Sprintf("%s/p2p/registry/%s/%s", impl.baseUrl, *otp, *token)

	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("X-Api-Key", impl.apiKey)
	response, err := client.Do(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer response.Body.Close()
	log.Debugln(response.StatusCode, url)
	return response.StatusCode == http.StatusAccepted
}

func (impl *BlitzshareApiImpl) RegisterAsPeer(multiAddr *string, otp, mode *string) *string {
	var token *string = nil
	payload, _ := json.Marshal(map[string]string{
		"multiAddr": *multiAddr,
		"otp":       *otp,
		"mode":      *mode,
	})
	url := fmt.Sprintf("%s/p2p/registry", impl.baseUrl)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("X-Api-Key", impl.apiKey)
	response, err := client.Do(req)
	defer response.Body.Close()
	req.Header.Set("X-Api-Key", impl.apiKey)
	if err != nil {
		log.Fatalln(err)
	}
	if response.StatusCode == http.StatusAccepted {
		ack := PeerRegistryAckResponse{}
		body, _ := ioutil.ReadAll(response.Body)
		err = json.Unmarshal(body, &ack)
		if err != nil {
			log.Fatalln(err)
		}
		log.Debugln("RegisterAsPeer", ack.Token, ack.AckId, url)
		token = &ack.Token
	}
	return token
}

func (impl *BlitzshareApiImpl) GetPeerConfig(oneTimePass *string) *P2pPeerRegistryResponse {
	var result *P2pPeerRegistryResponse = nil
	url := fmt.Sprintf("%s/p2p/registry/%s", impl.baseUrl, *oneTimePass)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Api-Key", impl.apiKey)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		peerAddress := P2pPeerRegistryResponse{}
		err = json.Unmarshal(body, &peerAddress)
		if err != nil {
			fmt.Println(err)
		}
		result = &peerAddress
	}
	return result
}

func (impl *BlitzshareApiImpl) GetBootstrapNode() *NodeConfigRespone {
	var result *NodeConfigRespone = nil
	url := fmt.Sprintf("%s/p2p/bootstrap-node", impl.baseUrl)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Api-Key", impl.apiKey)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		config := NodeConfigRespone{}
		err = json.Unmarshal(body, &config)
		if err == nil {
			result = &config
		}
	}
	return result
}
