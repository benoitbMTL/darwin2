package utils

import (
	"net/mail"
	"testing"
	"unicode"
)

func TestGeneratedSamplesCoverEmailAndPasswordPatterns(t *testing.T) {
	emailPatterns := make(map[string]struct{})
	passwordPatterns := make(map[string]struct{})

	for index := range 1482 {
		sample := generateFakeData(uint64(index))
		if _, err := mail.ParseAddress(BuildEmail(*sample)); err != nil {
			t.Fatalf("sample %d generated invalid email %q: %v", index, BuildEmail(*sample), err)
		}
		if index < 100 {
			emailPatterns[metricPattern(BuildEmail(*sample))] = struct{}{}
		}
		passwordPatterns[metricPattern(sample.Password)] = struct{}{}
	}

	// Some local/domain combinations collapse to the same C/N/S sequence
	// because FortiWeb also categorizes the @ separator as S.
	if got := len(emailPatterns); got != 81 {
		t.Fatalf("email pattern count = %d, want 81", got)
	}
	if got := len(passwordPatterns); got != 1482 {
		t.Fatalf("password pattern count = %d, want 1482", got)
	}
}

func TestGeneratedNamesContainOnlyRealisticNameCharacters(t *testing.T) {
	for index := range 500 {
		sample := generateFakeData(uint64(index))
		for field, value := range map[string]string{"first name": sample.FirstName, "last name": sample.LastName} {
			if value == "" {
				t.Fatalf("sample %d has an empty %s", index, field)
			}
			for _, character := range value {
				if !unicode.IsLetter(character) && character != ' ' && character != '-' && character != '\'' {
					t.Fatalf("sample %d %s %q contains invalid character %q", index, field, value, character)
				}
			}
		}
	}
}

func TestPasswordPatternCatalog(t *testing.T) {
	if got := len(passwordBlockPatterns); got != 1482 {
		t.Fatalf("password pattern catalog contains %d patterns, want 1482", got)
	}
	for _, pattern := range passwordBlockPatterns {
		if len(pattern) < 3 || len(pattern) > 9 {
			t.Fatalf("invalid metric length for pattern %q", pattern)
		}
		for index := 1; index < len(pattern); index++ {
			if pattern[index] == pattern[index-1] {
				t.Fatalf("pattern %q contains adjacent identical categories", pattern)
			}
		}
	}
}

func TestThreeThousandSamplesCoverEveryPasswordPattern(t *testing.T) {
	for _, start := range []int{0, 1, 731, 1481} {
		patterns := make(map[string]struct{})
		for offset := range 3000 {
			password := realisticPassword(start + offset)
			if len(password) < 8 || len(password) > 64 {
				t.Fatalf("password length = %d, want 8..64 for %q", len(password), password)
			}
			patterns[metricPattern(password)] = struct{}{}
		}
		if got := len(patterns); got != len(passwordBlockPatterns) {
			t.Fatalf("start %d covers %d password patterns, want %d", start, got, len(passwordBlockPatterns))
		}
	}
}

func metricPattern(value string) string {
	pattern := make([]rune, 0, len(value))
	var previous rune
	for _, character := range value {
		category := 'S'
		if unicode.IsLetter(character) {
			category = 'C'
		} else if unicode.IsDigit(character) {
			category = 'N'
		}
		if category != previous {
			pattern = append(pattern, category)
			previous = category
		}
	}
	return string(pattern)
}
