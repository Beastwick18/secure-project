package phone

import (
	"fmt"
	"secure/util"
)

func ValidPhone(p string) bool {
	// Match internal extension number
	if util.Match(`^\d{5}$`, p) {
		return true
	}

	no_zero3 := `([1-9][0-9]{0,2})`
	no_zero_both := `(([1-9][0-9]{0,2})|([1-9][0-9]{0,1}))`
	// no_zero2 := `[1-9][0-9]{0,1}`

	// Match 011 +1 (12) 123-1234 like numbers
	if util.Match(fmt.Sprintf(`^(\d{3}\s?)?(\+?%s\s?)?(\(%s\)\s?)?\d{3}\-\d{4}$`, no_zero3, no_zero_both), p) {
		return true
	}

	// Match any 123-123-1234 or 1-123-123-1234 with any valid seperator
	seps := []string{`\.`, `\-`, ` `}
	for _, sep := range seps {
		if util.Match(fmt.Sprintf(`^(\d{3}%s)?(\+?%s%s)?\d{2,3}%s\d{3}%s\d{4}$`, sep, no_zero3, sep, sep, sep), p) {
			return true
		}
	}

	// Danish AA AA AA AA
	seps = []string{`\.`, ` `}
	for _, sep := range seps {
		if util.Match(fmt.Sprintf(`^(\d{3}%s)?(\+?45%s)?\d{2}%s\d{2}%s\d{2}%s\d{2}$`, sep, sep, sep, sep, sep), p) {
			return true
		}
	}

	// Ten digit AAAAA AAAAA
	seps = []string{`\.`, ` `}
	for _, sep := range seps {
		if util.Match(fmt.Sprintf(`^\d{5}%s\d{5}$`, sep), p) {
			return true
		}
	}

	return false
}
