package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const BASE_CONFIG_API_URL = "http://0.0.0.0"

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
	log.Infoln("RegisterAsPeer", resp.Body, resp.Status)
	return resp.Status == "200 OK"
}
