package main

import (
	"flag"

	"github.com/sethgrid/simplequest/dungeon"
	"github.com/sethgrid/simplequest/quester"
	"github.com/sethgrid/simplequest/utils"
)

func main() {
	// get config, run game
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set flag to enable debug logs")
	flag.Parse()

	utils.Debug = debug

	d := dungeon.MakeDungeon()
	d.NewCell("outworld 1,1").
		Description("you are in a clearing").
		AddDestination("outworld 1,0", "north", "north, a dark forest").
		AddDestination("outworld 2,1", "east", "to the east, large spires peak over the hills")

	d.NewCell("outworld 1,0").
		Description("you find yourself in a dark forest").
		AddDestination("outworld 1,1", "south", "to the south, a clearing")

	d.NewCell("outworld 2,1").
		Description("you stand before a great castle. To the north is dense forest, and to the south, a great chasm").
		AddDestination("outworld 1.1", "west", "to the west lies a clearing").
		AddDestination("castle 0,1", "steps|door|rusty door|east", "Up the large steps stands a large door")

	d.NewCell("castle 0,1").
		Description("you arrive at the top of the steps, a little out of breath.").
		AddDestination("castle 1,1", "east|door|rusty door", "a giant rusty door, it is open just enough to fit through.").
		AddDestination("outworld 2,1", "steps|west", "going down the steps to the west leads away from the castle.")

	d.NewCell("castle 1,1").
		Description("you enter the geat entry way. Soon, there will be three doors.").
		AddDestination("castle 0,1", "rusty door", "the rusty door is behind you").
		AddDestination("castle 2,0", "blue door", "a blue door").
		AddDestination("castle 2,0", "red door", "a red door with a blue lock on it").
		AddDestination("castle 2,0", "green door", "a green door with a white lock on it")

	// blue door
	d.NewCell("castle 2,0").
		Description("in the center of room is a large stone table").
		AddDestination("castle 1,1", "door|blue door", "the door out is on the eastern wall")

	// red door
	d.NewCell("castle 2,1").
		Description("There is a white door. It is sealed shut with no apparent lock, inscribed on its doors is bit.ly/sddffs").
		AddDestination("castle 1,1", "door|red door", "the door out is on the eastern wall")

	// green door
	d.NewCell("castle 2,2").
		Description("How did you get in this room, there was no white key...").
		AddDestination("castle 1,1", "door|green door", "the door out is on the eastern wall")

	// white room
	d.NewCell("castle 3,2").
		Description("How did you get through to this room, the door was sealed...").
		AddDestination("castle 1,1", "door|white door", "the door out is on the eastern wall")

	quest := quester.NewQuest(quester.NewPlayer("foo"), d)
	quest.Start()
}

/* outworld

              [ forest 1,0        ] [ forest 2,0           ] [ forest 3,0 ]
[ ocean 0,1 ] [ starting zone 1,1 ] [ castle enterance 2,1 ] [ forest 3,1 ]
              [ river 1,2         ] [ chasm 2,2            ] [ chasm  3,2 ]
*/

/* castle
                                          [ blue door  2,0 ]
[ castle enterance 0,1]  [ main room 1,1] [ red door   2,1 ][ white room 3,2 ]
						                  [ green door 2,2 ]
*/
