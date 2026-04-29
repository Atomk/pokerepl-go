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
			input:    "  hello world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "    ",
			expected: []string{},
		},
		{
			input: `
			`,
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		expected := c.expected
		if len(actual) != len(expected) {
			t.Errorf("actual and expects have a different number of elements\nactual: %v\nexpected: %v", actual, expected)
			continue
		}
		for i := range actual {
			if actual[i] != expected[i] {
				t.Errorf("%s != %s", actual[i], expected[i])
			}
		}
	}
}
