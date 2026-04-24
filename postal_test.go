package postal_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/airtrafik/postal"
)

func TestPostal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Postal Suite")
}

var _ = Describe("Expand", Ordered, func() {
	var p *postal.Postal

	BeforeAll(func() {
		var err error
		p, err = postal.New()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterAll(func() {
		p.Close()
	})

	It("expands a simple English address", func() {
		expansions := p.Expand("123 Main St")
		Expect(expansions).To(ContainElement("123 main street"))
	})

	It("expands with English language options", func() {
		opts := postal.DefaultExpandOptions()
		opts.Languages = []string{"en"}

		expansions := p.ExpandWithOptions("30 West Twenty-sixth St Fl No. 7", opts)
		Expect(expansions).To(ContainElement("30 west 26th street floor number 7"))

		expansions = p.ExpandWithOptions("Thirty W 26th St Fl #7", opts)
		Expect(expansions).To(ContainElement("30 west 26th street floor number 7"))
	})

	It("expands with multilingual options", func() {
		opts := postal.DefaultExpandOptions()
		opts.Languages = []string{"en", "fr", "de"}

		expansions := p.ExpandWithOptions("st", opts)
		Expect(expansions).To(ContainElement("sankt"))
		Expect(expansions).To(ContainElement("saint"))
	})

	It("expands non-ASCII addresses", func() {
		expansions := p.Expand("Friedrichstraße 128, Berlin, Germany")
		Expect(expansions).To(ContainElement("friedrich strasse 128 berlin germany"))
	})
})
