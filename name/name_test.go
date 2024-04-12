package name

import (
	"testing"
)

func TestValidNames(t *testing.T) {
	valid_names := []string{
		// "O'Malley, John F.",
		"Bruce Schneier",
		"Bruce F. Schneier",
		"Schneier, Bruce",
		"Schneier, Bruce Wayne",
		"O'Malley, John F.",
		"John O'Malley-Smith",
		"Cher",
	}
	for _, n := range valid_names {
		if res := ValidName(n); !res {
			t.Fatalf(`Input: "%s", Result: "%v", Expected: "%v"`, n, res, true)
		}
	}
}

func TestInvalidNames(t *testing.T) {
	invalid_names := []string{
		"Ron O''Henry",
		"Ron O'Henry-Smith-Barnes",
		"L33t Hacker",
		"<Script>alert(“XSS”)</Script>",
		"Brad Everett Samuel Smith",
		"select * from users;",
	}

	for _, n := range invalid_names {
		if res := ValidName(n); res {
			t.Fatalf(`Input: "%s", Result: "%v", Expected: "%v"`, n, res, false)
		}
	}
}
