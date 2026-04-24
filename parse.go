package postal

/*
#cgo pkg-config: libpostal
#include <libpostal/libpostal.h>
#include <stdlib.h>
*/
import "C"

import (
	"unicode/utf8"
	"unsafe"
)

// ParseOptions controls address parsing behavior.
type ParseOptions struct {
	Language string
	Country  string
}

// ParsedComponent represents a single labeled component of a parsed address.
type ParsedComponent struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// Parse parses an address string into labeled components using default options.
func (p *Postal) Parse(address string) []ParsedComponent {
	return p.ParseWithOptions(address, ParseOptions{})
}

// ParseWithOptions parses an address string using the given options.
func (p *Postal) ParseWithOptions(address string, options ParseOptions) []ParsedComponent {
	if !utf8.ValidString(address) {
		return nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

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
