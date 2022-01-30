package random_test

import (
	"testing"

	"bootstrap.cli/app/services/random"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRandomService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "random test")
}

var _ = Describe("test random module", func() {
	var rnd random.Rnd
	BeforeSuite(func() {
		rnd = random.NewRnd()
	})
	Context("given random instance", func() {
		It("expected to generate random word sequence", func() {
			words := rnd.GenerateRandomWordSequence()
			Expect(words).To(Not(BeNil()))
			Expect(len(*words) > 6).To(BeTrue())
			Expect(words).To(Not(Equal(rnd.GenerateRandomWordSequence())))
		})
	})
})
