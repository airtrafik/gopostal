package postal

/*
#cgo pkg-config: libpostal
#include <libpostal/libpostal.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"os"
	"sync"
	"unicode/utf8"
	"unsafe"
)

var mu sync.Mutex

// Setup initializes libpostal and the address parser with the default data directory.
func Setup() error {
	return SetupDataDir("")
}

// SetupDataDir initializes libpostal and the address parser with a custom data directory.
// If dataDir is empty, it checks the LIBPOSTAL_DATA_DIR environment variable.
// If neither is set, libpostal's compiled-in default is used.
func SetupDataDir(dataDir string) error {
	if dataDir == "" {
		dataDir = os.Getenv("LIBPOSTAL_DATA_DIR")
	}

	mu.Lock()
	defer mu.Unlock()

	if dataDir != "" {
		cDataDir := C.CString(dataDir)
		defer C.free(unsafe.Pointer(cDataDir))
		if !bool(C.libpostal_setup_datadir(cDataDir)) {
			return fmt.Errorf("libpostal setup failed for data directory: %s", dataDir)
		}
		if !bool(C.libpostal_setup_parser_datadir(cDataDir)) {
			return fmt.Errorf("libpostal parser setup failed for data directory: %s", dataDir)
		}
	} else {
		if !bool(C.libpostal_setup()) {
			return fmt.Errorf("libpostal setup failed: data directory may be missing")
		}
		if !bool(C.libpostal_setup_parser()) {
			return fmt.Errorf("libpostal parser setup failed")
		}
	}
	return nil
}

// Teardown frees libpostal resources. Call this when done using the parser package.
func Teardown() {
	mu.Lock()
	defer mu.Unlock()
	C.libpostal_teardown()
	C.libpostal_teardown_parser()
}

type ParserOptions struct {
	Language string
	Country  string
}

func getDefaultParserOptions() ParserOptions {
	return ParserOptions{
		Language: "",
		Country:  "",
	}
}

var parserDefaultOptions = getDefaultParserOptions()

type ParsedComponent struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func ParseAddressOptions(address string, options ParserOptions) []ParsedComponent {
	if !utf8.ValidString(address) {
		return nil
	}

	mu.Lock()
	defer mu.Unlock()

	cAddress := C.CString(address)
	defer C.free(unsafe.Pointer(cAddress))

	cOptions := C.libpostal_get_address_parser_default_options()
	if options.Language != "" {
		cLanguage := C.CString(options.Language)
		defer C.free(unsafe.Pointer(cLanguage))

		cOptions.language = cLanguage
	}

	if options.Country != "" {
		cCountry := C.CString(options.Country)
		defer C.free(unsafe.Pointer(cCountry))

		cOptions.country = cCountry
	}

	cAddressParserResponsePtr := C.libpostal_parse_address(cAddress, cOptions)

	if cAddressParserResponsePtr == nil {
		return nil
	}

	cAddressParserResponse := *cAddressParserResponsePtr

	cNumComponents := cAddressParserResponse.num_components
	cComponents := cAddressParserResponse.components
	cLabels := cAddressParserResponse.labels

	numComponents := uint64(cNumComponents)

	parsedComponents := make([]ParsedComponent, numComponents)

	// Accessing a C array
	cComponentsPtr := (*[1 << 30](*C.char))(unsafe.Pointer(cComponents))[:numComponents:numComponents]
	cLabelsPtr := (*[1 << 30](*C.char))(unsafe.Pointer(cLabels))[:numComponents:numComponents]

	var i uint64
	for i = 0; i < numComponents; i++ {
		parsedComponents[i] = ParsedComponent{
			Label: C.GoString(cLabelsPtr[i]),
			Value: C.GoString(cComponentsPtr[i]),
		}
	}

	C.libpostal_address_parser_response_destroy(cAddressParserResponsePtr)

	return parsedComponents
}

func ParseAddress(address string) []ParsedComponent {
	return ParseAddressOptions(address, parserDefaultOptions)
}
