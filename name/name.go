package name

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

func ValidName(name string) bool {
	first := `(([A-Z][a-z]*'[A-Z]?[a-z]+)|([A-Z][a-z]+))`
	last := `(([A-Z][a-z]*'[A-Z]?[a-z]+)|([A-Z][a-z]+))`
	last = `(` + last + `(-` + last + `)?)` // Allow optional -secondLastName
	middle := `(([A-Z]\.?)|([A-Z][a-z]*'[A-Z]?[a-z]+)|([A-Z][a-z]+))`
	if match(`^`+first+`( `+middle+`)?( `+last+`)?$`, name) {
		return true
	}
	if match(`^`+last+`\, `+first+`( `+middle+`)?$`, name) {
		return true
	}
	return false
}
