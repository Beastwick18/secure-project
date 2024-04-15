package name

import (
	"fmt"
	"secure/util"
	"strings"
)

func ValidName(name string) bool {
	first := `([a-z]+('[a-z]+)?)`
	last := `([a-z]+('[a-z]+)?)`
	last = fmt.Sprintf(`(%s(-%s)?)`, last, last) // Allow optional -secondLastName
	middle := `(([a-z]\.?)|([a-z]+('[a-z]+)?))`
	name = strings.ToLower(name)

	fml := fmt.Sprintf(`^%s( %s)?( %s)?$`, first, middle, last)
	lfm := fmt.Sprintf(`^%s\, %s( %s)?`, last, first, middle)
	return util.Match(fml, name) || util.Match(lfm, name)
}
