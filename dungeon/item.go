package dungeon

import (
	"github.com/sethgrid/simplequest/parser"
)

// Item represents anything in a room that can be altered
type Item struct {
	Name        string
	Takable     bool
	InInventory bool
	Movable     bool
	InRoomDesc  string
	InInvDesc   string
	// Action takes in an verb and an optional list of items.
	// the list can be one or two items.
	// If two, it is expected to be item 1 as the initiating item and
	// item 2 as the receiving item.
	// Eg "smash", "table", "rock"
	Action func(command parser.Parsed, inventory ...[]*Item) string
}

// NewItem is a helper to initialize an Item
func NewItem(name string, takable bool, movable bool, inRoomDesc string, inInvDesc string, action func(command parser.Parsed, inventory ...[]*Item) string) *Item {
	return &Item{
		Name:       name,
		Takable:    takable,
		Movable:    movable,
		InRoomDesc: inRoomDesc,
		InInvDesc:  inInvDesc,
		Action:     action,
	}
}
