package blitzshare_test

import (
	"fmt"
	"net/http"
	"testing"

	"bootstrap.cli/app/config"
	"bootstrap.cli/app/services/blitzshare"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRandomService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "blitzshare api test")
}

var _ = Describe("test blitzshare api module", func() {
	var mockedConfig config.AppConfig
	var api blitzshare.BlitzshareApi
	BeforeSuite(func() {
		mockedConfig = config.AppConfig{
			P2pBoostrapNodeIp:   "",
			P2pBoostrapNodeId:   "",
			P2pBoostrapNodePort: 0,
			BlitzshareApiUrl:    "http://test.host",
			LocalP2pPeerIp:      "",
		}
		api = blitzshare.NewBlitzsahreApi(&mockedConfig)
		httpmock.Activate()
	})
	Context("GetBootstrapNode Tests", func() {
		It("expected node config to be nil for http status 200 (StatusOK)", func() {
			nodeResponse := blitzshare.NodeConfigRespone{
				NodeId: "node-test-id",
				Port:   1234,
			}
			resp, _ := httpmock.NewJsonResponder(http.StatusOK, nodeResponse)
			url := fmt.Sprintf("%s/p2p/bootstrap-node", mockedConfig.BlitzshareApiUrl)
			httpmock.RegisterResponder("GET", url, resp)
			nodeConf := api.GetBootstrapNode()
			Expect(nodeConf.NodeId).To(Equal(nodeResponse.NodeId))
			Expect(nodeConf.Port).To(Equal(nodeResponse.Port))
		})
		It("expected node config to be nil for http status 500 (StatusInternalServerError)", func() {
			resp, _ := httpmock.NewJsonResponder(http.StatusInternalServerError, nil)
			url := fmt.Sprintf("%s/p2p/bootstrap-node", mockedConfig.BlitzshareApiUrl)
			httpmock.RegisterResponder("GET", url, resp)
			nodeConf := api.GetBootstrapNode()
			Expect(nodeConf).To(BeNil())
		})
	})
	Context("GetPeerConfig Tests", func() {
		It("expected peer config for 200 (StatusOK)", func() {
			otp := "otp-otp-otp"
			nodeResponse := blitzshare.P2pPeerRegistryResponse{
				MultiAddr: "test-MultiAddr",
				Otp:       otp,
				Mode:      "chat",
			}
			resp, _ := httpmock.NewJsonResponder(http.StatusOK, nodeResponse)
			url := fmt.Sprintf("%s/p2p/registry/%s", mockedConfig.BlitzshareApiUrl, otp)
			httpmock.RegisterResponder("GET", url, resp)
			peer := api.GetPeerConfig(&otp)
			Expect(peer.Mode).To(Equal(nodeResponse.Mode))
			Expect(peer.Otp).To(Equal(nodeResponse.Otp))
			Expect(peer.Mode).To(Equal(nodeResponse.Mode))
		})
		It("expected peer config a nil for http status 500 (StatusInternalServerError)", func() {
			otp := "otp-otp-otp"
			resp, _ := httpmock.NewJsonResponder(http.StatusInternalServerError, nil)
			url := fmt.Sprintf("%s/p2p/registry/%s", mockedConfig.BlitzshareApiUrl, otp)
			httpmock.RegisterResponder("GET", url, resp)
			peer := api.GetPeerConfig(&otp)
			Expect(peer).To(BeNil())

		})
	})
})
