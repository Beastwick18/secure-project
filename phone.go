package main

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

func valid_phone(p string) bool {
	// Subscriber only 123-1234
	// if match, _ := regexp.MatchString("^[0-9]{3}([\\.\\-]|\\s*)[0-9]{4}$", p); match {
	// 	return true
	// }
	// Match (234)123-1234 or 123-1234
	if match(`^(1\s?)?(\(\d{3}\)\s?)?\d{3}\-\d{4}$`, p) {
		return true
	}

	// Match any 123-123-1234 or 1-123-123-1234 with any valid seperator
	seps := []string{`\.`, `\-`, ` `, ``}
	for _, sep := range seps {
		if match(`^(1`+sep+`)?\d{3}`+sep+`\d{3}`+sep+`\d{4}$`, p) {
			return true
		}
	}
	return false
	// match, err := regexp.MatchString("^(\\+?[0-9]{1,3})-?([0-9]{3})$", p)
	// return match, err
}
