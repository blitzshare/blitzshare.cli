package stream_test

import (
	"bufio"
	"os"
	"testing"

	"bootstrap.cli/app/services/stream"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStrService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Stream test")
}

var _ = Describe("Test stream module", func() {
	Context("WriteStreamFromStdin", func() {
		It("expected to callback to be called once stream had been writen to", func() {
			callback := make(chan string, 1)
			rw := bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdin))
			stream.WriteStreamFromStdin(rw, func() {
				callback <- "done"
			})
			rw.WriteString("helo from test")
			<-callback
			Expect(true).To(BeTrue())
		})
	})
})
