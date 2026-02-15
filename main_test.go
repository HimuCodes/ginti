package main_test

import (
	"testing"

	counter "github.com/HimuCodes/counter"
)

func TestCountWords(t *testing.T) {
	type testCase struct {
		name  string
		input string
		wants int
	}

	testCases := []struct {
		name  string
		input string
		wants int
	}{
		{
			name:  "5 words",
			input: "one two three four five",
			wants: 5,
		},
		{
			name:  "empty input",
			input: "",
			wants: 0,
		},
		{
			name:  "single space",
			input: " ",
			wants: 0,
		},
		{
			name:  "new lines",
			input: "one two three\nfour five",
			wants: 5,
		},
		{
			name:  "multiple spaces",
			input: "This is a sentence. This is another",
			wants: 7,
		},
		{
			name:  "Prefixed multiple spaces",
			input: "  Hello",
			wants: 1,
		},
		{
			name:  "Sufficed  multiple spaces",
			input: "Hello    ",
			wants: 1,
		},
		{
			name:  "Tab character in code",
			input: "Hello\tWord\n",
			wants: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := counter.CountWords([]byte(tc.input))
			if result != tc.wants {
				t.Logf("expected: %d got: %d", tc.wants, result)
				t.Fail()
			}
		})
	}
}
