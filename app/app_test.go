package app_test

import (
	"bufio"
	"os"
	"testing"

	"bootstrap.cli/app"
	"bootstrap.cli/app/config"
	"bootstrap.cli/app/dependencies"
	"bootstrap.cli/app/services/blitzshare"
	"bootstrap.cli/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
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
	Context("given app module", func() {
		It("expected StartPeer to return otp", func() {
			api := &mocks.BlitzshareApi{}
			api.On("RegisterAsPeer",
				mock.AnythingOfType("string"),
				mock.MatchedBy(func(input interface{}) bool {
					return true
				})).Return(true)
			p2p := &mocks.P2p{}
			p2p.On("StartPeer", mock.MatchedBy(func(input interface{}) bool {
				return true
			}), mock.MatchedBy(func(input interface{}) bool {
				return true
			}), mock.MatchedBy(func(input interface{}) bool {
				return true
			})).Return("tcp://0.0.0.0/whatever")
			rnd := &mocks.Rnd{}
			otp := "clogwood-bristle-overwrap-benzdifuran"
			rnd.On("GenerateRandomWordSequence").Return(&otp)
			clipboard := &mocks.ClipBoard{}
			clipboard.On("CopyToClipBoard", mock.MatchedBy(func(input interface{}) bool {
				return true
			})).Return()
			dep := &dependencies.Dependencies{
				Config:        &mockedConfig,
				BlitzshareApi: api,
				P2p:           p2p,
				Rnd:           rnd,
				ClipBoard:     clipboard,
			}
			peerOTP := app.StartPeer(dep)
			Expect(otp).To(Equal(*peerOTP))
		})
	})
	Context("given ConnectToPeerOTP", func() {
		It("expected ConnectToPeerOTP to connect", func() {
			api := &mocks.BlitzshareApi{}
			apiResponse := &blitzshare.PeerAddress{
				MultiAddr: "tcp://0.0.0.0/whatever",
			}
			api.On("GetPeerAddr", mock.MatchedBy(func(input interface{}) bool {
				return true
			})).Return(apiResponse)
			p2p := &mocks.P2p{}

			rw := bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdin))
			p2p.On("ConnectToPeer", mock.MatchedBy(func(input interface{}) bool {
				return true
			}), mock.MatchedBy(func(input interface{}) bool {
				return true
			}), mock.MatchedBy(func(input interface{}) bool {
				return true
			})).Return(rw)

			otp := "clogwood-bristle-overwrap-benzdifuran"
			dep := &dependencies.Dependencies{
				Config:        &mockedConfig,
				BlitzshareApi: api,
				P2p:           p2p,
			}
			address := app.ConnectToPeerOTP(dep, &otp)
			Expect(address).To(Equal(apiResponse.MultiAddr))
		})
	})
})
