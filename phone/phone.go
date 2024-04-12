package phone

import (
	"log"
	"regexp"
)

func match(pattern string, input string) bool {
	match, err := regexp.MatchString(pattern, input)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	return match
}

func ValidPhone(p string) bool {
	// Match internal extension number
	if match(`^\d{5}$`, p) {
		return true
	}

	no_zero3 := `([1-9][0-9]{0,2})`
	no_zero_both := `(([1-9][0-9]{0,2})|([1-9][0-9]{0,1}))`
	// no_zero2 := `[1-9][0-9]{0,1}`

	if match(`^(\d{3}\s?)?(\+?`+no_zero3+`\s?)?(\(`+no_zero_both+`\)\s?)?\d{3}\-\d{4}$`, p) {
		return true
	}

	// Match 1(234)123-1234, (234)123-1234, or 123-1234
	// if match(`^(\d{3}\s?)?(\+?`+no_zero3+`\s?)?(\((`+no_zero3+`)|(`+no_zero2+`)\)\s?)?\d{3}\-\d{4}$`, p) {
	// 	return true
	// }

	// Match any 123-123-1234 or 1-123-123-1234 with any valid seperator
	seps := []string{`\.`, `\-`, ` `}
	for _, sep := range seps {
		if match(`^(\d{3}`+sep+`)?(\+?`+no_zero3+sep+`)?\d{2,3}`+sep+`\d{3}`+sep+`\d{4}$`, p) {
			return true
		}
	}

	// Danish AA AA AA AA
	seps = []string{`\.`, ` `}
	for _, sep := range seps {
		if match(`^(\d{3}`+sep+`)?(\+?45`+sep+`)?\d{2}`+sep+`\d{2}`+sep+`\d{2}`+sep+`\d{2}$`, p) {
			return true
		}
	}

	// Ten digit AAAAA AAAAA
	seps = []string{`\.`, ` `}
	for _, sep := range seps {
		if match(`^\d{5}`+sep+`\d{5}$`, p) {
			return true
		}
	}

	return false
}
