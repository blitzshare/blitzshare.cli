package config_test

import (
	"os"
	"testing"

	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	setUp()
	cfg, err := config.Load()

	assert.Nil(t, err, "Unable to log the config")
	assert.Equal(t, "12D3KooWFsVgymCfHTvhECN3wmwE7kFmYTNg1n7F2CZ3iMqwt2TP", cfg.P2pBoostrapNodeId)
	assert.Equal(t, "10.100.212.158", cfg.P2pBoostrapNodeIp)

	tearDown()
}

func setUp() {
	_ = os.Setenv("P2P_BOOTSTRAP_NODE_IP", "10.100.212.158")
	_ = os.Setenv("P2P_BOOTSTRAP_NODE_ID", "12D3KooWFsVgymCfHTvhECN3wmwE7kFmYTNg1n7F2CZ3iMqwt2TP")
}

func tearDown() {
	_ = os.Unsetenv("P2P_BOOTSTRAP_NODE_IP")
	_ = os.Unsetenv("P2P_BOOTSTRAP_NODE_ID")
}
