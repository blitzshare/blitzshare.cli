package config

import "github.com/kelseyhightower/envconfig"

type AppConfig struct {
	P2pBoostrapNodeIp   string `envconfig:"P2P_BOOTSTRAP_NODE_IP"`
	P2pBoostrapNodeId   string `envconfig:"P2P_BOOTSTRAP_NODE_ID"`
	P2pBoostrapNodePort int    `envconfig:"PORT" default:"63785"`
	BlitzshareApiUrl    string `envconfig:"BLITZSHARE_API_URL"`
	LocalP2pPeerIp      string
}

func Load() (*AppConfig, error) {
	LoadEnvironment()
	cfg := AppConfig{}
	cfg.LocalP2pPeerIp = "0.0.0.0"
	err := envconfig.Process("", &cfg)
	return &cfg, err
}
