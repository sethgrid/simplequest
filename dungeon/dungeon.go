package dungeon

import (
	"strings"

	"github.com/sethgrid/simplequest/utils"
)

var NotTakable = false
var Takable = true
var NoInventoryDesc = ""
var NoAction = func(verb string, item ...*Item) string { return "" }
var Locked = true
var Open = true

// Dungeon contains unexported properties to facilitate creating a dungeon.
type Dungeon struct {
	startingCell string           // cellID
	m            map[string]*Cell // key is cellID
}

// StartingCell returns the enterence to this dungeon
func (d *Dungeon) StartingCell() string {
	return d.startingCell
}

// LoadCell takes a cellID and returns that cell and true. If that cell is not available,
// it returns the starting cell and false to signify that the intended cell was unavailable.
func (d *Dungeon) LoadCell(cellID string) (*Cell, bool) {
	if cell, ok := d.m[cellID]; ok {
		return cell, ok
	}
	return d.m[d.startingCell], false
}

// Cell represents a given poition in a dungeon / map.
type Cell struct {
	ID               string
	description      string
	destinations     map[string]destination
	dungeon          Dungeon
	items            []*Item
	doors            map[string]*Door // door names -> door
	doorDestinations map[*Door]string // door -> cell id
}

type destination struct {
	description  string // to the north it is dark
	name         string // north
	cellID       string // some id, maybe '11,3'
	hiddenOption bool   // don't print it when listing options (used for hidding redundant options like door vs large door)
}

// MakeDungeon initializes an empty dungeon.
func MakeDungeon() *Dungeon {
	return &Dungeon{m: make(map[string]*Cell)}
}

// NewCell creates a new cell. The first call to this method creates the initial dungeon starting cell.
func (d *Dungeon) NewCell(id string) *Cell {
	if len(d.m) == 0 {
		d.startingCell = id
	}
	cell := &Cell{
		ID:               id,
		destinations:     make(map[string]destination),
		doors:            make(map[string]*Door),
		doorDestinations: make(map[*Door]string),
	}
	d.m[id] = cell
	return cell
}

// Description is the base description provided to the user. Destinations (and eventually items and mosters) augment the base description.
func (c *Cell) Description(s string) *Cell {
	c.description = s
	return c
}

// AddDestination creates a path out of this current cell.
// Different paths need different calls to AddDestination.
// A single path can have multiple "names" buy splitting with pipe.
// Ex: cell.AddDestination("castle3 0,0", "east|door|large door|enterance", "you see the enterance, a large door to the east")
func (c *Cell) AddDestination(cellID, name, localDescription string) *Cell {
	variations := strings.Split(name, "|")
	for i, variation := range variations {
		hiddenOption := false
		if i > 0 {
			hiddenOption = true
		}
		c.destinations[variation] = destination{cellID: cellID, description: localDescription, hiddenOption: hiddenOption}
	}
	return c
}

// AddItem to the cell
func (c *Cell) AddItem(item *Item) *Cell {
	c.items = append(c.items, item)
	return c
}

// AddDoor to the cell
func (c *Cell) AddDoor(cellID string, name string, door *Door) *Cell {
	variations := strings.Split(name, "|")
	for i, variation := range variations {
		c.doors[variation] = door
		if i > 0 {
			c.doors[variation] = door
		} else {
			c.doors[variation] = door
			c.doorDestinations[door] = cellID
		}
	}
	return c
}

// GetDoor will look for a door of the given name in this cell and return it if it exists
func (c *Cell) GetDoor(name string) (*Door, bool) {
	for _, door := range c.doors {
		if door.Name == name {
			return door, true
		}
	}
	return nil, false
}

// GetItem will look for an item of the given name in this cell and return it if it exists
func (c *Cell) GetItem(name string) (*Item, bool) {
	for _, item := range c.items {
		if item.Name == name {
			return item, true
		}
	}
	return nil, false
}

// RemoveItem disassociates the item with the cell
func (c *Cell) RemoveItem(name string) {
	itemIndex := -1
	for i, item := range c.items {
		if item.Name == name {
			itemIndex = i
			break
		}
	}
	if itemIndex >= 0 {
		// remove the item from the room's item list
		c.items = append(c.items[:itemIndex], c.items[itemIndex+1:]...)
	}
}

// GetDestinationID ...
func (c *Cell) GetDestinationID(name string) (string, bool) {
	utils.Debugf("looking for destination %q", name)
	if destination, ok := c.destinations[name]; ok {
		return destination.cellID, true
	}
	if door, ok := c.doors[name]; ok {
		return c.doorDestinations[door], true
	}
	return "", false
}

// Items retreives all items from the cell
func (c *Cell) Items() []*Item {
	return c.items
}

// Prompt combines the description with destinations for this cell (and eventually items and monsters) and prepars the output the user will see.
func (c *Cell) Prompt(promptChar string) string {
	utils.Debugf("current cell: %s", c.ID)
	prompt := c.description + "\n"

	for _, destination := range c.destinations {
		if destination.hiddenOption {
			continue
		}
		prompt += destination.description + "\n"
	}

	for door := range c.doorDestinations {
		prompt += door.Description + "\n"
	}

	for _, item := range c.items {
		if !item.Hidden {
			prompt += item.InRoomDesc + "\n"
		}
	}

	if promptChar != "" {
		prompt += promptChar
	}

	return prompt
}
