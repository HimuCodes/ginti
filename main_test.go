package main

import "testing"

func TestCountWords(t *testing.T) {
	input := "one two three four five"

	wants := 5

	result := CountWords([]byte(input))

	if result != wants {
		t.Log("expected:", wants, "got:", result)
		t.Fail()
	}

	input = ""
	wants = 0

	result = CountWords([]byte(input))

	if result != wants {
		t.Log("expected:", wants, "got:", result)
		t.Fail()
	}

	input = " "
	wants = 0

	result = CountWords([]byte(input))
	if result != wants {
		t.Log("expected:", wants, "got:", result)
		t.Fail()
	}
}
