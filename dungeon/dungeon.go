package dungeon

type Dungeon struct {
	m map[string]*Cell // key is cellID
}

type Cell struct {
	ID           string
	description  string
	destinations []destination
	dungeon      Dungeon
}

type destination struct {
	description string
	cellID      string
}

/* outer world

              [ forest 1,0        ] [ forest 2,0           ] [ forest 3,0 ]
[ ocean 0,1 ] [ starting zone 1,1 ] [ castle enterance 2,2 ] [ forest 3,1 ]
              [ river 1,2         ] [ chasm 2,3            ] [ chasm  3,2 ]
*/

/* castle
                                          [ blue door  2,0 ]
[ castle enterance 0,1]  [ main room 1,1] [ red door   2,1 ][ light room 3,2 ]
						                  [ green door 2,2 ]
*/

func MakeDungeon() *Dungeon {
	return &Dungeon{m: make(map[string]*Cell)}
}

func (d *Dungeon) NewCell(id string) *Cell {
	return &Cell{ID: id}
}

func (c *Cell) Description(s string) *Cell {
	c.description = s
	return c
}

func (c *Cell) AddDestination(cellID, localDescription string) {
	c.destinations = append(c.destinations, destination{cellID: cellID, description: localDescription})
}
