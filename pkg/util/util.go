package util

import (
	"log"
	"regexp"
)

func Match(pattern string, input string) bool {
	match, err := regexp.MatchString(pattern, input)
	if err != nil {
		log.Printf("Failed to parse regex:\n%s", err)
		return false
	}
	return match
}

func Contains[S ~[]E, E comparable](s S, v E) bool {
	for i := range s {
		if v == s[i] {
			return true
		}
	}
	return false
}
