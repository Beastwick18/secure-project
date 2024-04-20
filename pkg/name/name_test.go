package name

import (
	"testing"
)

func TestValidNames(t *testing.T) {
	valid_names := []string{
		"Bruce Schneier",
		"Bruce F. Schneier",
		"Schneier, Bruce",
		"Schneier, Bruce Wayne",
		"O'Malley, John F.",
		"John O'Malley-Smith",
		"Cher",

		// Student-provided
		"Steven B. Culwell",
		"Steven Bradley Culwell",
		"Steven Culwell",
		"Stèvèn Cülwèll",
		"Culwell-O'Culwell, Steven B.",
		"Раз два три",
		"Раз два три'три-три",
	}
	for _, n := range valid_names {
		t.Run(n, func(t *testing.T) {
			if res := ValidName(n); !res {
				t.Errorf(`Input: "%s", Result: "%v", Expected: "%v"`, n, res, true)
			}
		})
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

		// Student-provided
		"B. Steven, Culwell",
		"Culwell-O'Culwell, B. Steven",
		"Steven Culwell-Culwell-Culwell",
		"B. Culwell-O'Culwell Steven",
		"Раз два три'три-три два",
	}

	for _, n := range invalid_names {
		t.Run(n, func(t *testing.T) {
			if res := ValidName(n); res {
				t.Errorf(`Input: "%s", Result: "%v", Expected: "%v"`, n, res, false)
			}
		})
	}
}
