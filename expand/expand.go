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

// Setup initializes libpostal and the language classifier with the default data directory.
func Setup() error {
	return SetupDataDir("")
}

// SetupDataDir initializes libpostal and the language classifier with a custom data directory.
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
		if !bool(C.libpostal_setup_language_classifier_datadir(cDataDir)) {
			return fmt.Errorf("libpostal language classifier setup failed for data directory: %s", dataDir)
		}
	} else {
		if !bool(C.libpostal_setup()) {
			return fmt.Errorf("libpostal setup failed: data directory may be missing")
		}
		if !bool(C.libpostal_setup_language_classifier()) {
			return fmt.Errorf("libpostal language classifier setup failed")
		}
	}
	return nil
}

// Teardown frees libpostal resources. Call this when done using the expand package.
func Teardown() {
	mu.Lock()
	defer mu.Unlock()
	C.libpostal_teardown()
	C.libpostal_teardown_language_classifier()
}

const (
	AddressNone        = C.LIBPOSTAL_ADDRESS_NONE
	AddressAny         = C.LIBPOSTAL_ADDRESS_ANY
	AddressName        = C.LIBPOSTAL_ADDRESS_NAME
	AddressHouseNumber = C.LIBPOSTAL_ADDRESS_HOUSE_NUMBER
	AddressStreet      = C.LIBPOSTAL_ADDRESS_STREET
	AddressUnit        = C.LIBPOSTAL_ADDRESS_UNIT
	AddressLevel       = C.LIBPOSTAL_ADDRESS_LEVEL
	AddressStaircase   = C.LIBPOSTAL_ADDRESS_STAIRCASE
	AddressEntrance    = C.LIBPOSTAL_ADDRESS_ENTRANCE
	AddressCategory    = C.LIBPOSTAL_ADDRESS_CATEGORY
	AddressNear        = C.LIBPOSTAL_ADDRESS_NEAR
	AddressToponym     = C.LIBPOSTAL_ADDRESS_TOPONYM
	AddressPostalCode  = C.LIBPOSTAL_ADDRESS_POSTAL_CODE
	AddressPoBox       = C.LIBPOSTAL_ADDRESS_PO_BOX
	AddressAll         = C.LIBPOSTAL_ADDRESS_ALL
)

type ExpandOptions struct {
	Languages              []string
	AddressComponents      uint16
	LatinAscii             bool
	Transliterate          bool
	StripAccents           bool
	Decompose              bool
	Lowercase              bool
	TrimString             bool
	ReplaceWordHyphens     bool
	DeleteWordHyphens      bool
	ReplaceNumericHyphens  bool
	DeleteNumericHyphens   bool
	SplitAlphaFromNumeric  bool
	DeleteFinalPeriods     bool
	DeleteAcronymPeriods   bool
	DropEnglishPossessives bool
	DeleteApostrophes      bool
	ExpandNumex            bool
	RomanNumerals          bool
}

var cDefaultOptions = C.libpostal_get_default_options()

func GetDefaultExpansionOptions() ExpandOptions {
	return ExpandOptions{
		Languages:              nil,
		AddressComponents:      uint16(cDefaultOptions.address_components),
		LatinAscii:             bool(cDefaultOptions.latin_ascii),
		Transliterate:          bool(cDefaultOptions.transliterate),
		StripAccents:           bool(cDefaultOptions.strip_accents),
		Decompose:              bool(cDefaultOptions.decompose),
		Lowercase:              bool(cDefaultOptions.lowercase),
		TrimString:             bool(cDefaultOptions.trim_string),
		ReplaceWordHyphens:     bool(cDefaultOptions.replace_word_hyphens),
		DeleteWordHyphens:      bool(cDefaultOptions.delete_word_hyphens),
		ReplaceNumericHyphens:  bool(cDefaultOptions.replace_numeric_hyphens),
		DeleteNumericHyphens:   bool(cDefaultOptions.delete_numeric_hyphens),
		SplitAlphaFromNumeric:  bool(cDefaultOptions.split_alpha_from_numeric),
		DeleteFinalPeriods:     bool(cDefaultOptions.delete_final_periods),
		DeleteAcronymPeriods:   bool(cDefaultOptions.delete_acronym_periods),
		DropEnglishPossessives: bool(cDefaultOptions.drop_english_possessives),
		DeleteApostrophes:      bool(cDefaultOptions.delete_apostrophes),
		ExpandNumex:            bool(cDefaultOptions.expand_numex),
		RomanNumerals:          bool(cDefaultOptions.roman_numerals),
	}
}

var libpostalDefaultOptions = GetDefaultExpansionOptions()

func ExpandAddressOptions(address string, options ExpandOptions) []string {
	if !utf8.ValidString(address) {
		return nil
	}

	mu.Lock()
	defer mu.Unlock()

	cAddress := C.CString(address)
	defer C.free(unsafe.Pointer(cAddress))

	var char_ptr *C.char
	ptr_size := unsafe.Sizeof(char_ptr)

	cOptions := C.libpostal_get_default_options()
	if options.Languages != nil {
		cLanguages := C.calloc(C.size_t(len(options.Languages)), C.size_t(ptr_size))
		cLanguagesPtr := (*[1 << 30](*C.char))(unsafe.Pointer(cLanguages))

		defer C.free(unsafe.Pointer(cLanguages))

		for i := 0; i < len(options.Languages); i++ {
			cLang := C.CString(options.Languages[i])
			defer C.free(unsafe.Pointer(cLang))
			cLanguagesPtr[i] = cLang
		}

		cOptions.languages = (**C.char)(cLanguages)
		cOptions.num_languages = C.size_t(len(options.Languages))
	} else {
		cOptions.num_languages = 0
	}

	cOptions.address_components = C.uint16_t(options.AddressComponents)
	cOptions.latin_ascii = C.bool(options.LatinAscii)
	cOptions.transliterate = C.bool(options.Transliterate)
	cOptions.strip_accents = C.bool(options.StripAccents)
	cOptions.decompose = C.bool(options.Decompose)
	cOptions.lowercase = C.bool(options.Lowercase)
	cOptions.trim_string = C.bool(options.TrimString)
	cOptions.replace_word_hyphens = C.bool(options.ReplaceWordHyphens)
	cOptions.delete_word_hyphens = C.bool(options.DeleteWordHyphens)
	cOptions.replace_numeric_hyphens = C.bool(options.ReplaceNumericHyphens)
	cOptions.delete_numeric_hyphens = C.bool(options.DeleteNumericHyphens)
	cOptions.split_alpha_from_numeric = C.bool(options.SplitAlphaFromNumeric)
	cOptions.delete_final_periods = C.bool(options.DeleteFinalPeriods)
	cOptions.delete_acronym_periods = C.bool(options.DeleteAcronymPeriods)
	cOptions.drop_english_possessives = C.bool(options.DropEnglishPossessives)
	cOptions.delete_apostrophes = C.bool(options.DeleteApostrophes)
	cOptions.expand_numex = C.bool(options.ExpandNumex)
	cOptions.roman_numerals = C.bool(options.RomanNumerals)

	var cNumExpansions = C.size_t(0)

	cExpansions := C.libpostal_expand_address(cAddress, cOptions, &cNumExpansions)

	numExpansions := uint64(cNumExpansions)

	var expansions = make([]string, numExpansions)

	// Accessing a C array
	cExpansionsPtr := (*[1 << 30](*C.char))(unsafe.Pointer(cExpansions))

	var i uint64
	for i = 0; i < numExpansions; i++ {
		expansions[i] = C.GoString(cExpansionsPtr[i])
	}

	C.libpostal_expansion_array_destroy(cExpansions, cNumExpansions)
	return expansions
}

func ExpandAddress(address string) []string {
	return ExpandAddressOptions(address, libpostalDefaultOptions)
}
