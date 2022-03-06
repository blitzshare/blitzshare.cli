package config_test

import (
	"os"
	"testing"

	"bootstrap.cli/app/config"
	"github.com/stretchr/testify/assert"
)

const (
	BoostrapNodeIp   = "11.111.111.111"
	BlitzshareApiUrl = "22.22.22.22"
	BlitzshareApiKey = "test-api-key"
	BoostrapNodePort = 63785
)

func TestConfig(t *testing.T) {
	setUp()
	cfg, err := config.Load()

	assert.Nil(t, err, "Unable to log the config")
	assert.Equal(t, BoostrapNodeIp, cfg.P2pBoostrapNodeIp)
	assert.Equal(t, BlitzshareApiUrl, cfg.BlitzshareApiUrl)
	assert.Equal(t, BlitzshareApiKey, cfg.BlitzshareApiKey)
	assert.Equal(t, BoostrapNodePort, BoostrapNodePort)

	tearDown()
}

func setUp() {
	_ = os.Setenv("BLITZSHARE_API_URL", BlitzshareApiUrl)
	_ = os.Setenv("BLITZSHARE_API_KEY", BlitzshareApiKey)
	_ = os.Setenv("P2P_BOOTSTRAP_NODE_IP", BoostrapNodeIp)
}

func tearDown() {
	_ = os.Unsetenv("P2P_BOOTSTRAP_NODE_IP")
	_ = os.Unsetenv("BLITZSHARE_API_KEY")
	_ = os.Unsetenv("BLITZSHARE_API_URL")
}
