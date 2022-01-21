package config_test

import (
	"os"
	"testing"

	"bootstrap.cli/app/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	setUp()
	cfg, err := config.Load()

	assert.Nil(t, err, "Unable to log the config")
	assert.Equal(t, "11.111.111.111", cfg.P2pBoostrapNodeIp)
	assert.Equal(t, "22.22.22.22", cfg.BlitzshareApiUrl)

	tearDown()
}

func setUp() {
	_ = os.Setenv("BLITZSHARE_API_URL", "22.22.22.22")
	_ = os.Setenv("P2P_BOOTSTRAP_NODE_IP", "11.111.111.111")

}

func tearDown() {
	_ = os.Unsetenv("P2P_BOOTSTRAP_NODE_IP")
}
