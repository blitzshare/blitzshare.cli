package app_test

import (
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/config"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/dependencies"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/services/blitzshare"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestApp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry test")
}

var _ = Describe("App tests", func() {
	BeforeSuite(func() {

	})
	Context("wha", func() {
		It("what what", func() {
			api := &mocks.BlitzshareApi{}
			api.On("RegisterAsPeer").Return(true)
			p2p := &mocks.P2p{}
			p2p.On("StartPeer").Return(nil)
			c := &config.AppConfig{
				P2pBoostrapNodeIp:   "",
				P2pBoostrapNodeId:   "",
				P2pBoostrapNodePort: 0,
				BlitzshareApiUrl:    "",
				LocalP2pPeerIp:      "",
			}
			dep := &dependencies.Dependencies{
				Config:        c,
				BlitzshareApi: blitzshare.NewBlitzsahreApi(c),
				P2p:           nil,
			}
			app.StartPeer(dep)
			Expect(1).To(Equal(2))
		})
	})
})
