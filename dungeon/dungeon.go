package dungeon

import (
	"strings"

	"github.com/sethgrid/simplequest/utils"
)

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
	ID           string
	description  string
	destinations map[string]destination
	dungeon      Dungeon
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
	cell := &Cell{ID: id, destinations: make(map[string]destination)}
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

// GetDestinationID ...
func (c *Cell) GetDestinationID(name string) (string, bool) {
	if destination, ok := c.destinations[name]; ok {
		return destination.cellID, true
	}
	return "", false
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

	if promptChar != "" {
		prompt += promptChar
	}

	return prompt
}
