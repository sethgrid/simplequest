package dungeon

import "github.com/sethgrid/simplequest/parser"

// Door ...
type Door struct {
	Name        string
	Description string
	IsOpen      bool
	IsLocked    bool
	// PerformActionAndPrompt takes in an verb and an optional list of items.
	// the list can be one or two items.
	// If two, it is expected to be item 1 as the initiating item and
	// item 2 as the receiving item.
	// Eg "unlock", "blue door", "blue key"
	PerformActionAndPrompt func(cell *Cell, command parser.Parsed, inventory ...*Item) string
}

// NewDoor is a helper to initialize a Door
func NewDoor(name string, isLocked bool, isOpen bool, description string, action func(cell *Cell, command parser.Parsed, inventory ...*Item) string) *Door {
	return &Door{
		Name:                   name,
		IsLocked:               isLocked,
		IsOpen:                 isOpen,
		Description:            description,
		PerformActionAndPrompt: action,
	}
}
