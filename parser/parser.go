package parser

import (
	"strings"

	"github.com/sethgrid/simplequest/utils"
)

type Parsed struct {
	Sentence string
	Action   string

	Object     string
	Identifier string
}

var skipWords = []string{"an", "the", "up", "a", "to", "around", "through", "over", "beside", "on"}

func Parse(sentence string) Parsed {
	parsed := &Parsed{Sentence: sentence}
	words := strings.Split(sentence, " ")
	var gotFirstWord bool
	var gotSecondWord bool
	for _, word := range words {
		word = strings.ToLower(word)
		if utils.StrInList(word, skipWords) {
			continue
		}
		if word == "and" {
			utils.Debugf("compound commands not supported yet")
			break
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
	parsed.Action = verbSynonym(parsed.Action)
	return *parsed
}

var goVerbs = []string{"go", "run", "walk", "travel", "head", "venture", "approach"}
var lookVerbs = []string{"look", "inspect", "read"}
var takeVerbs = []string{"take", "steal", "get"}
var useVerbs = []string{"use", "activate"}
var exitVerbs = []string{"exit", "quit", "leave"}

// verbSynonym hones down the verbs we have to explicitly handle
func verbSynonym(someVerb string) string {
	switch {
	case utils.StrInList(someVerb, goVerbs):
		return "go"
	case utils.StrInList(someVerb, lookVerbs):
		return "look"
	case utils.StrInList(someVerb, takeVerbs):
		return "take"
	case utils.StrInList(someVerb, useVerbs):
		return "use"
	case utils.StrInList(someVerb, exitVerbs):
		return "exit"
	default:
		return ""
	}
}
