package dependencies

import (
	cfg "github.com/blitzshare/blitzshare.bootstrap.client.cli/app/config"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/services/blitzshare"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/services/p2p"
)

type Dependencies struct {
	Config        *cfg.AppConfig
	BlitzshareApi blitzshare.BlitzshareApi
	P2p           p2p.P2p
}

func NewDependencies(config *cfg.AppConfig) (*Dependencies, error) {
	return &Dependencies{
			Config:        config,
			BlitzshareApi: blitzshare.NewBlitzsahreApi(config),
			P2p:           p2p.NewP2p(),
		},
		nil
}
