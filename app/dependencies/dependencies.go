package dependencies

import (
	cfg "github.com/blitzshare/blitzshare.bootstrap.client.cli/app/config"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/services/blitzshare"
)

type Dependencies struct {
	Config        *cfg.AppConfig
	BlitzshareApi *blitzshare.BlitzshareApiImpl
}

func NewDependencies(config *cfg.AppConfig) (*Dependencies, error) {
	return &Dependencies{Config: config, BlitzshareApi: blitzshare.New(config)}, nil
}
