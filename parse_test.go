package postal_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/airtrafik/postal"
)

var _ = Describe("Parse", func() {
	It("parses a US address into components", func() {
		components := p.Parse("781 Franklin Ave Crown Heights Brooklyn NYC NY 11216 USA")

		Expect(components).To(Equal([]postal.ParsedComponent{
			{Label: "house_number", Value: "781"},
			{Label: "road", Value: "franklin ave"},
			{Label: "suburb", Value: "crown heights"},
			{Label: "city_district", Value: "brooklyn"},
			{Label: "city", Value: "nyc"},
			{Label: "state", Value: "ny"},
			{Label: "postcode", Value: "11216"},
			{Label: "country", Value: "usa"},
		}))
	})

	It("marshals to JSON correctly", func() {
		components := p.Parse("781 Franklin Ave Crown Heights Brooklyn NYC NY 11216 USA")

		data, err := json.Marshal(components)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(Equal(
			`[{"label":"house_number","value":"781"},{"label":"road","value":"franklin ave"},{"label":"suburb","value":"crown heights"},{"label":"city_district","value":"brooklyn"},{"label":"city","value":"nyc"},{"label":"state","value":"ny"},{"label":"postcode","value":"11216"},{"label":"country","value":"usa"}]`,
		))
	})

	It("parses with language and country options", func() {
		components := p.ParseWithOptions("781 Franklin Ave Brooklyn NY 11216", postal.ParseOptions{
			Language: "en",
			Country:  "us",
		})
		Expect(components).NotTo(BeEmpty())
		Expect(components[0].Label).To(Equal("house_number"))
		Expect(components[0].Value).To(Equal("781"))
	})

	It("round-trips through JSON", func() {
		components := p.Parse("781 Franklin Ave Crown Heights Brooklyn NYC NY 11216 USA")

		data, err := json.Marshal(components)
		Expect(err).NotTo(HaveOccurred())

		var unmarshaled []postal.ParsedComponent
		Expect(json.Unmarshal(data, &unmarshaled)).To(Succeed())
		Expect(unmarshaled).To(Equal(components))
	})
})
