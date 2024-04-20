package name

import (
	"fmt"
	"secure/pkg/util"
	"strings"
)

func ValidName(name string) bool {
	first := `([\p{L}]+('[\p{L}]+)?)`
	last := `([\p{L}]+('[\p{L}]+)?)`
	last = fmt.Sprintf(`(%s(-%s)?)`, last, last) // Allow optional -secondLastName
	middle := `(([\p{L}]\.?)|([\p{L}]+('[\p{L}]+)?))`
	name = strings.ToLower(name)

	fml := fmt.Sprintf(`^%s( %s)?( %s)?$`, first, middle, last)
	lfm := fmt.Sprintf(`^%s\, %s( %s)?$`, last, first, middle)
	return util.Match(fml, name) || util.Match(lfm, name)
}
