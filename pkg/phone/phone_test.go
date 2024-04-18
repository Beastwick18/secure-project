package phone

import (
	"testing"
)

func TestValidPhones(t *testing.T) {
	valid_phones := []string{
		"12345",
		"(703)111-2121",
		"123-1234",
		"+1(703)111-2121",
		"1-703-111-2121",
		"+12-12-123-1234",
		"+32 (21) 212-2324",
		"1(703)123-1234",
		"011 701 111 1234",
		"12345.12345",
		"011 1 703 111 1234",
	}

	for _, p := range valid_phones {
		match := ValidPhone(p)
		if !match {
			t.Fatalf(`Input: "%s", Result: "%v", Expected: "%v"`, p, match, true)
		}
	}
}

func TestInvalidPhones(t *testing.T) {
	invalid_phones := []string{
		"123",
		"1/703/123/1234",
		"Nr 102-123-1234",
		"<script>alert(“XSS”)</script>",
		"7031111234",
		"+1234 (201) 123-1234",
		"(001) 123-1234",
		"(01) 123-1234",
		"+01 (703) 123-1234",
		"(703) 123-1234 ext 204",
	}
	for _, p := range invalid_phones {
		match := ValidPhone(p)
		if match {
			t.Fatalf(`Input: "%s", Result: "%v", Expected: "%v"`, p, match, false)
		}
	}
}
