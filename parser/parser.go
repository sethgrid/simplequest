package parser

import (
	"strings"

	"github.com/sethgrid/simplequest/utils"
)

// Parsed represents a command that has been filtered for use in the quest
type Parsed struct {
	Sentence  string
	RawAction string
	Action    string

	Object     string
	Identifier string

	ActionObject     string
	ActionIdentifier string
}

var skipWords = []string{"an", "the", "up", "a", "to", "around", "through", "over", "beside", "on", "in", "is", "my"}

// Parse takes in a simple sentence and transforms it into a parsed structure that can be used to perform an action.
func Parse(sentence string) Parsed {
	parsed := &Parsed{Sentence: sentence}

	// strip some common punctuation
	sentence = strings.Replace(sentence, ".", "", -1)
	sentence = strings.Replace(sentence, "?", "", -1)
	sentence = strings.Replace(sentence, "!", "", -1)

	words := strings.Split(sentence, " ")
	var gotFirstWord bool
	var gotSecondWord bool

	//open the door with the key
	for _, word := range words {
		if word == "" {
			continue
		}
		word = strings.ToLower(word)
		if utils.StrInList(word, skipWords) {
			continue
		}
		if word == "and" {
			utils.Debugf("compound commands not supported yet")
			break
		}
		if word == "with" {
			// set up parser to receive a new object
			parsed.ActionObject = parsed.Object
			parsed.ActionIdentifier = parsed.Identifier
			gotSecondWord = false
			continue
		}
		if !gotFirstWord {
			parsed.RawAction = word
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

		// The last thing we can parse is the Action identifier. If we have that, stop scanning the sentence.
		if parsed.ActionIdentifier != "" {
			break
		}
	}
	if parsed.ActionObject != "" {
		// if we had a sentance with two objects, they probably want the second object to be the action object
		parsed.Object, parsed.ActionObject = parsed.ActionObject, parsed.Object
		parsed.Identifier, parsed.ActionIdentifier = parsed.ActionIdentifier, parsed.Identifier
	}
	parsed.Action = verbSynonym(parsed.RawAction)
	return *parsed
}

var goVerbs = []string{"go", "run", "walk", "travel", "head", "venture", "approach", "climb", "enter"}
var lookVerbs = []string{"look", "inspect", "read", "what", "what's", "whats"}
var takeVerbs = []string{"take", "steal", "get", "grab", "remove"}
var useVerbs = []string{"use", "activate"}
var sayVerbs = []string{"say", "speak", "tell", "yell", "whisper", "shout"}
var exitVerbs = []string{"exit", "quit"}

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
	case utils.StrInList(someVerb, sayVerbs):
		return "say"
	default:
		return someVerb
	}
}
