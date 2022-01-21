package dependencies

import (
	cfg "bootstrap.cli/app/config"
	"bootstrap.cli/app/services/blitzshare"
	"bootstrap.cli/app/services/clipboard"
	"bootstrap.cli/app/services/p2p"
	rnd "bootstrap.cli/app/services/random"
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
