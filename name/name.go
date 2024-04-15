package name

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

func match(pattern string, input string) bool {
	match, err := regexp.MatchString(pattern, input)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	return match
}

func ValidName(name string) bool {
	first := `([a-z]+('[a-z]+)?)`
	last := `([a-z]+('[a-z]+)?)`
	last = fmt.Sprintf(`(%s(-%s)?)`, last, last) // Allow optional -secondLastName
	middle := `(([a-z]\.?)|([a-z]+('[a-z]+)?))`
	name = strings.ToLower(name)

	fml := fmt.Sprintf(`^%s( %s)?( %s)?$`, first, middle, last)
	lfm := fmt.Sprintf(`^%s\, %s( %s)?`, last, first, middle)
	return match(fml, name) || match(lfm, name)
}
