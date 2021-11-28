package str_test

import (
	"testing"

	"github.com/blitzshare/blitzshare.bootstrap.client.cli/app/services/str"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStrService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry test")
}

var _ = Describe("Test str modue", func() {
	Context("Given a registry", func() {
		It("expected to return the same string if no chars needs to be removed", func() {
			strVal := "/ip4/10.101.18.26/tcp/63785/p2p/12D3KooWPGR"
			Expect(str.SanatizeStr(strVal)).To(Equal(strVal))
		})
		It("expected to remove new line", func() {
			Expect(str.SanatizeStr("12D3KooWPGR\n")).To(Equal("12D3KooWPGR"))
			Expect(str.SanatizeStr("12D3KooWPGR\n\n")).To(Equal("12D3KooWPGR"))
			Expect(str.SanatizeStr("\n12D3KooWPGR\n\n")).To(Equal("12D3KooWPGR"))
		})
		It("expected to remove white spaces", func() {
			Expect(str.SanatizeStr("12 D3K ooWPGR\n")).To(Equal("12D3KooWPGR"))
			Expect(str.SanatizeStr("12 D3K ooWPGR ")).To(Equal("12D3KooWPGR"))
			Expect(str.SanatizeStr(" 12 D3K ooWPGR ")).To(Equal("12D3KooWPGR"))
		})
	})
})
