package dependencies

import (
	cfg "github.com/blitzshare/blitzshare.bootstrap.client.cli/app/config"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/services/blitzshare"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/services/clipboard"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/services/p2p"
	rnd "github.com/blitzshare/blitzshare.bootstrap.client.cli/app/services/random"
)

type Dependencies struct {
	Config        *cfg.AppConfig
	BlitzshareApi blitzshare.BlitzshareApi
	P2p           p2p.P2p
	Rnd           rnd.Rnd
	ClipBoard     clipboard.ClipBoard
}

func NewDependencies(config *cfg.AppConfig) (*Dependencies, error) {
	return &Dependencies{
			Config:        config,
			BlitzshareApi: blitzshare.NewBlitzsahreApi(config),
			P2p:           p2p.NewP2p(),
			Rnd:           rnd.NewRnd(),
			ClipBoard:     clipboard.NewClipBoard(),
		},
		nil
}
