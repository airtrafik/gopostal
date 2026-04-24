package postal_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/airtrafik/postal"
)

func TestPostal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Postal Suite")
}

// findDataDir locates the libpostal data directory.
// Checks LIBPOSTAL_DATA_DIR env var, then common platform paths.
func findDataDir() string {
	if dir := os.Getenv("LIBPOSTAL_DATA_DIR"); dir != "" {
		return dir
	}

	// macOS: Homebrew
	if out, err := exec.Command("brew", "--prefix", "libpostal").Output(); err == nil {
		dir := filepath.Join(strings.TrimSpace(string(out)), "share", "libpostal")
		if hasLibpostalData(dir) {
			return dir
		}
	}

	// Alpine: apk package
	if hasLibpostalData("/usr/share/libpostal") {
		return "/usr/share/libpostal"
	}

	return ""
}

// hasLibpostalData checks if a directory contains libpostal data files.
func hasLibpostalData(dir string) bool {
	// Check for transliteration.dat — present in both Homebrew and Alpine packages.
	_, err := os.Stat(filepath.Join(dir, "transliteration", "transliteration.dat"))
	return err == nil
}

var _ = Describe("New", func() {
	It("returns an error for an invalid data directory", func() {
		_, err := postal.New(postal.WithDataDir("/nonexistent/path"))
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("/nonexistent/path"))
	})
})

var _ = Describe("Expand", Ordered, func() {
	var p *postal.Postal

	BeforeAll(func() {
		dataDir := findDataDir()
		Expect(dataDir).NotTo(BeEmpty(), "libpostal data directory not found — install libpostal with data files")

		var err error
		p, err = postal.New(postal.WithDataDir(dataDir))
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

	It("returns nil for invalid UTF-8 input", func() {
		expansions := p.Expand("\xff\xfe")
		Expect(expansions).To(BeNil())
	})
})
