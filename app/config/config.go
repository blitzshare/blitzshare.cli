package config

import "github.com/kelseyhightower/envconfig"

type AppConfig struct {
	Settings Settings
}

type Settings struct {
	P2pBoostrapNodeIp   string `envconfig:"P2P_BOOTSTRAP_NODE_IP"`
	P2pBoostrapNodeId   string `envconfig:"P2P_BOOTSTRAP_NODE_ID"`
	P2pBoostrapNodePort int    `envconfig:"PORT" default:"63785"`
	BlitzshareApiUrl    string `envconfig:"BLITZSHARE_API_URL"`
}

func Load() (*AppConfig, error) {
	LoadEnvironment()
	cfg := AppConfig{}
	err := envconfig.Process("settings", &cfg)
	return &cfg, err
}
