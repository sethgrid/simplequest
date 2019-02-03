package parser

import "strings"

type Parsed struct {
	Sentence string
	Action   string

	Object     string
	Identifier string
}

var skipWords = []string{"an", "the", "a", "to", "around", "through", "over", "beside", "on"}

func Parse(sentence string) Parsed {
	parsed := &Parsed{Sentence: sentence}
	words := strings.Split(sentence, " ")
	var gotFirstWord bool
	var gotSecondWord bool
	for _, word := range words {
		word = strings.ToLower(word)
		if inList(word, skipWords) {
			continue
		}
		if !gotFirstWord {
			parsed.Action = word
			gotFirstWord = true
			continue
		}
		if !gotSecondWord {
			parsed.Object = word
			gotSecondWord = true
			continue
		}
		// already have an action and an object, we must have missed an identifyer
		// ie, we thought it was open the door, but it is open the blue door
		// move blue from object to identifer and put door as object
		parsed.Identifier = parsed.Object
		parsed.Object = word
		break
	}
	return *parsed
}

func inList(needle string, haystack []string) bool {
	for _, element := range haystack {
		if needle == element {
			return true
		}
	}
	return false
}
