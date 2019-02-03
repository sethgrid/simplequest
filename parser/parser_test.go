package parser

import "testing"

func TestSimpleSentences(t *testing.T) {
	tests := []struct {
		in     string
		parsed Parsed
	}{
		{
			in: "go north",
			parsed: Parsed{
				Sentence:   "go north",
				Action:     "go",
				Object:     "north",
				Identifier: "",
			},
		},
		{
			in: "open the door",
			parsed: Parsed{
				Sentence:   "open the door",
				Action:     "open",
				Object:     "door",
				Identifier: "",
			},
		},
		{
			in: "open the blue door",
			parsed: Parsed{
				Sentence:   "open the blue door",
				Action:     "open",
				Object:     "door",
				Identifier: "blue",
			},
		},
		{
			in: "open",
			parsed: Parsed{
				Sentence:   "open",
				Action:     "open",
				Object:     "",
				Identifier: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			out := Parse(test.in)
			if !expectedParse(out, test.parsed) {
				t.Errorf("\ngot:\n%#v\nwant:\n%#v", out, test.parsed)
			}
		})
	}
}

func expectedParse(got, want Parsed) bool {
	if got.Sentence != want.Sentence {
		return false
	}
	if got.Action != want.Action {
		return false
	}
	if got.Identifier != want.Identifier {
		return false
	}
	if got.Object != want.Object {
		return false
	}
	return true
}
