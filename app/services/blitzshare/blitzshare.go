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
	GetPeerConfig(oneTimePass *string) *P2pPeerRegistryResponse
	GetBootstrapNode() *NodeConfigRespone
	DeregisterAsPeer(otp, token *string) bool
}

type BlitzshareApiImpl struct {
	BaseUrl string
	ApiKey  string
}

func New(config *cfg.AppConfig) BlitzshareApi {
	return &BlitzshareApiImpl{
		BaseUrl: config.BlitzshareApiUrl,
		ApiKey:  config.BlitzshareApiKey,
	}
}

func (impl *BlitzshareApiImpl) DeregisterAsPeer(otp, token *string) bool {
	url := fmt.Sprintf("%s/p2p/registry/%s/%s", impl.BaseUrl, *otp, *token)
	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("x-api-key", impl.ApiKey)
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
	client := &http.Client{}
	payload, _ := json.Marshal(map[string]string{
		"multiAddr": *multiAddr,
		"otp":       *otp,
		"mode":      *mode,
	})
	url := fmt.Sprintf("%s/p2p/registry", impl.BaseUrl)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Add("x-api-key", impl.ApiKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusAccepted {
		ack := PeerRegistryAckResponse{}
		body, _ := ioutil.ReadAll(resp.Body)
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
	client := &http.Client{}
	url := fmt.Sprintf("%s/p2p/registry/%s", impl.BaseUrl, *oneTimePass)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("x-api-key", impl.ApiKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
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
	client := &http.Client{}
	url := fmt.Sprintf("%s/p2p/bootstrap-node", impl.BaseUrl)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("x-api-key", impl.ApiKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		config := NodeConfigRespone{}
		err = json.Unmarshal(body, &config)
		if err == nil {
			result = &config
		}
	}
	return result
}
