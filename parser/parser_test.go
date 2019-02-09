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
				RawAction:  "go",
				Object:     "north",
				Identifier: "",
			},
		},
		{
			in: "open the door",
			parsed: Parsed{
				Sentence:   "open the door",
				Action:     "open",
				RawAction:  "open",
				Object:     "door",
				Identifier: "",
			},
		},
		{
			in: "open the blue door",
			parsed: Parsed{
				Sentence:   "open the blue door",
				Action:     "open",
				RawAction:  "open",
				Object:     "door",
				Identifier: "blue",
			},
		},
		{
			in: "open",
			parsed: Parsed{
				Sentence:   "open",
				Action:     "open",
				RawAction:  "open",
				Object:     "",
				Identifier: "",
			},
		},
		{
			in: "travel north",
			parsed: Parsed{
				Sentence:   "travel north",
				Action:     "go",
				RawAction:  "travel",
				Object:     "north",
				Identifier: "",
			},
		},
		{
			in: "open the door with the key",
			parsed: Parsed{
				Sentence:         "open the door with the key",
				Action:           "open",
				RawAction:        "open",
				Object:           "door",
				Identifier:       "",
				ActionObject:     "key",
				ActionIdentifier: "",
			},
		},
		{
			in: "unlock the red lock with the blue key",
			parsed: Parsed{
				Sentence:         "unlock the red lock with the blue key",
				Action:           "unlock",
				RawAction:        "unlock",
				Object:           "lock",
				Identifier:       "red",
				ActionObject:     "key",
				ActionIdentifier: "blue",
			},
		},
		{
			in: "unlock the red lock with the blue key",
			parsed: Parsed{
				Sentence:         "unlock the red lock with the blue key",
				Action:           "unlock",
				RawAction:        "unlock",
				Object:           "lock",
				Identifier:       "red",
				ActionObject:     "key",
				ActionIdentifier: "blue",
			},
		},
		{
			in: "unlock the lock with the master key",
			parsed: Parsed{
				Sentence:         "unlock the lock with the master key",
				Action:           "unlock",
				RawAction:        "unlock",
				Object:           "lock",
				Identifier:       "",
				ActionObject:     "key",
				ActionIdentifier: "master",
			},
		},
		{
			in: "inventory",
			parsed: Parsed{
				Sentence:         "inventory",
				Action:           "inventory",
				RawAction:        "inventory",
				Object:           "",
				Identifier:       "",
				ActionObject:     "",
				ActionIdentifier: "",
			},
		},
		{
			in: "what's in my inventory",
			parsed: Parsed{
				Sentence:         "what's in my inventory",
				Action:           "look",
				RawAction:        "what's",
				Object:           "inventory",
				Identifier:       "",
				ActionObject:     "",
				ActionIdentifier: "",
			},
		},
		////////////// not supported yet ////////////////
		// {
		// 	in: "use the key to unlock the lock",
		// 	parsed: Parsed{
		// 		Sentence:         "use the key to unlock the lock",
		// 		Action:           "use",
		// 		RawAction:        "use",
		// 		Object:           "lock",
		// 		Identifier:       "",
		// 		ActionObject:     "key",
		// 		ActionIdentifier: "",
		// 	},
		// },
		// {
		// 	in: "use the key to on the lock",
		// 	parsed: Parsed{
		// 		Sentence:         "use the key to on the lock",
		// 		Action:           "use",
		// 		RawAction:        "use",
		// 		Object:           "lock",
		// 		Identifier:       "",
		// 		ActionObject:     "key",
		// 		ActionIdentifier: "",
		// 	},
		// },
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
	if got.RawAction != want.RawAction {
		return false
	}
	if got.ActionIdentifier != want.ActionIdentifier {
		return false
	}
	if got.ActionObject != want.ActionObject {
		return false
	}
	return true
}
