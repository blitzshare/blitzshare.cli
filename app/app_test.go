package app_test

import (
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/config"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/dependencies"
	"github.com/blitzshare/blitzshare.bootstrap.client.cli/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestApp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry test")
}

var _ = Describe("App tests", func() {
	var mockedConfig config.AppConfig

	BeforeSuite(func() {
		mockedConfig = config.AppConfig{
			P2pBoostrapNodeIp:   "",
			P2pBoostrapNodeId:   "",
			P2pBoostrapNodePort: 0,
			BlitzshareApiUrl:    "",
			LocalP2pPeerIp:      "",
		}
	})
	Context("wha", func() {
		It("what what", func() {
			api := &mocks.BlitzshareApi{}
			api.On("")
			api.On("RegisterAsPeer",
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string")).
				Return(func(multiAddr string, oneTimePass string) bool { return true })

			p2p := &mocks.P2p{}
			p2p.On("StartPeer", mock.MatchedBy(func(input interface{}) bool {
				return true
			}), mock.MatchedBy(func(input interface{}) bool {
				return true
			}), mock.MatchedBy(func(input interface{}) bool {
				return true
			})).Return(nil, "sd")
			dep := &dependencies.Dependencies{
				Config:        &mockedConfig,
				BlitzshareApi: api,
				P2p:           p2p,
			}
			h := app.StartPeer(dep)
			Expect(h).To(BeNil())
		})
	})
})
