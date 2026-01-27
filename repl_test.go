package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "HELLO WORLD",
			expected: []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice against the expected slice
		if len(c.expected) != len(actual) {
			t.Errorf("expected length and actual length don't match.\nExpected: %d\nActual:%d", len(c.expected), len(actual))
			t.Fail()
		}
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			if word != expectedWord {
				t.Errorf("words do not match.\nWord: %v\nExpected Word: %v", word, expectedWord)
				t.Fail()
			}
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
		}
	}
}
