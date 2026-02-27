package main_test

import (
	"bytes"
	"strings"
	"testing"

	counter "github.com/HimuCodes/counter"
)

func TestCountWords(t *testing.T) {
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
			r := strings.NewReader(tc.input)
			result := counter.CountWords(r)
			if result != tc.wants {
				t.Logf("expected: %d got: %d", tc.wants, result)
				t.Fail()
			}
		})
	}
}

func TestCountLines(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		wants int
	}{
		{
			name:  "Simple five words, 1 new line",
			input: "one two three four five\n",
			wants: 1,
		},
		{
			name:  "empty file",
			input: "",
			wants: 0,
		},
		{
			name:  "no new lines",
			input: "one two three four five",
			wants: 0,
		},
		{
			name:  "no new line at end",
			input: "one two three four five\n six",
			wants: 1,
		},
		{
			name:  "multi newline string",
			input: "\n\n\n\n",
			wants: 4,
		},
		{
			name:  "multi word line string",
			input: "one\ntwo\nthree\nfour\nfive\n",
			wants: 5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.input)

			res := counter.CountLines(r)

			if res != tc.wants {
				t.Logf("expected: %d got: %d", tc.wants, res)
				t.Fail()
			}
		})
	}
}

func TestCountBytes(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		wants int
	}{
		{
			name:  "five words",
			input: "one two three four five",
			wants: 23,
		},
		{
			name:  "empty file",
			input: "",
			wants: 0,
		},
		{
			name:  "all spaces",
			input: "       ",
			wants: 7,
		},
		{
			name:  "newlines and words",
			input: "one\ntwo\nthree\nfour",
			wants: 18,
		},
		{
			name:  "multibyte characters",
			input: "こんにちは",
			wants: 15,
		},
		{
			name:  "mixed ascii and unicode",
			input: "hello 世界",
			wants: 12,
		},
		{
			name:  "control characters and tabs",
			input: "col1\tcol2\r\n",
			wants: 11,
		},
		{
			name:  "null bytes",
			input: "\x00\x00\x00",
			wants: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.input)

			res := counter.CountBytes(r)

			if res != tc.wants {
				t.Logf("expected: %d got: %d", tc.wants, res)
				t.Fail()
			}
		})
	}
}

func TestGetCounts(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		wants counter.Counts
	}{
		{
			name:  "simple five words",
			input: "one two three four five\n",
			wants: counter.Counts{
				Lines: 1,
				Words: 5,
				Bytes: 24,
			},
		},
		{
			name:  "empty input",
			input: "",
			wants: counter.Counts{
				Lines: 0,
				Words: 0,
				Bytes: 0,
			},
		},
		{
			name:  "single space",
			input: " ",
			wants: counter.Counts{
				Lines: 0,
				Words: 0,
				Bytes: 1,
			},
		},
		{
			name:  "new lines in words",
			input: "one two three\nfour five",
			wants: counter.Counts{
				Lines: 1,
				Words: 5,
				Bytes: 23,
			},
		},
		{
			name:  "multiple spaces between words",
			input: "This is a sentence. This is another",
			wants: counter.Counts{
				Lines: 0,
				Words: 7,
				Bytes: 35,
			},
		},
		{
			name:  "prefixed multiple spaces",
			input: "  Hello",
			wants: counter.Counts{
				Lines: 0,
				Words: 1,
				Bytes: 7,
			},
		},
		{
			name:  "suffixed multiple spaces",
			input: "Hello    ",
			wants: counter.Counts{
				Lines: 0,
				Words: 1,
				Bytes: 9,
			},
		},
		{
			name:  "tab character in code",
			input: "Hello\tWord\n",
			wants: counter.Counts{
				Lines: 1,
				Words: 2,
				Bytes: 11,
			},
		},
		{
			name:  "no new lines",
			input: "one two three four five",
			wants: counter.Counts{
				Lines: 0,
				Words: 5,
				Bytes: 23,
			},
		},
		{
			name:  "no new line at end",
			input: "one two three four five\n six",
			wants: counter.Counts{
				Lines: 1,
				Words: 6,
				Bytes: 28,
			},
		},
		{
			name:  "multi newline string",
			input: "\n\n\n\n",
			wants: counter.Counts{
				Lines: 4,
				Words: 0,
				Bytes: 4,
			},
		},
		{
			name:  "multi word line string",
			input: "one\ntwo\nthree\nfour\nfive\n",
			wants: counter.Counts{
				Lines: 5,
				Words: 5,
				Bytes: 24,
			},
		},
		{
			name:  "all spaces",
			input: "       ",
			wants: counter.Counts{
				Lines: 0,
				Words: 0,
				Bytes: 7,
			},
		},
		{
			name:  "newlines and words",
			input: "one\ntwo\nthree\nfour",
			wants: counter.Counts{
				Lines: 3,
				Words: 4,
				Bytes: 18,
			},
		},
		{
			name:  "multibyte characters",
			input: "こんにちは",
			wants: counter.Counts{
				Lines: 0,
				Words: 1,
				Bytes: 15,
			},
		},
		{
			name:  "mixed ascii and unicode",
			input: "hello 世界",
			wants: counter.Counts{
				Lines: 0,
				Words: 2,
				Bytes: 12,
			},
		},
		{
			name:  "control characters and tabs",
			input: "col1\tcol2\r\n",
			wants: counter.Counts{
				Lines: 1,
				Words: 2,
				Bytes: 11,
			},
		},
		{
			name:  "null bytes",
			input: "\x00\x00\x00",
			wants: counter.Counts{
				Lines: 0,
				Words: 1,
				Bytes: 3,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.input)

			res := counter.GetCounts(r)

			if res != tc.wants {
				t.Logf("expected: %v got: %v", tc.wants, res)
				t.Fail()
			}
		})
	}
}

func TestPrint(t *testing.T) {
	testCases := []struct {
		name  string
		input struct {
			counts   counter.Counts
			filename string
		}
		wants string
	}{
		{
			name: "simple five words with filename",
			input: struct {
				counts   counter.Counts
				filename string
			}{
				counts: counter.Counts{
					Lines: 1,
					Words: 5,
					Bytes: 24,
				},
				filename: "words.txt",
			},
			wants: "       1        5       24 words.txt\n",
		},
		{
			name: "empty counts no filename",
			input: struct {
				counts   counter.Counts
				filename string
			}{
				counts: counter.Counts{
					Lines: 0,
					Words: 0,
					Bytes: 0,
				},
				filename: "",
			},
			wants: "       0        0        0\n",
		},
		{
			name: "large numbers with filename",
			input: struct {
				counts   counter.Counts
				filename string
			}{
				counts: counter.Counts{
					Lines: 1000,
					Words: 5000,
					Bytes: 30000,
				},
				filename: "big.txt",
			},
			wants: "    1000     5000    30000 big.txt\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buffer := &bytes.Buffer{}

			tc.input.counts.Print(buffer, tc.input.filename)

			if buffer.String() != tc.wants {
				t.Logf("expected: %q got: %q", tc.wants, buffer.String())
				t.Fail()
			}
		})
	}
}
