package quester

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sethgrid/simplequest/dungeon"
	"github.com/sethgrid/simplequest/parser"
	"github.com/sethgrid/simplequest/utils"
)

type Player struct {
	id        string
	inventory []*dungeon.Item
}

// NewPlayer initializes a player
func NewPlayer(id string) *Player {
	return &Player{
		id: id,
	}
}

// AddInventory adds to the players inventory
func (p *Player) AddInventory(item *dungeon.Item) {
	p.inventory = append(p.inventory, item)
}

// DescribeInventory does what you think it would
func (p *Player) DescribeInventory() string {
	if len(p.inventory) == 0 {
		return "you have nothing in your inventory"
	}
	description := "You have the following in your inventory:"
	for _, item := range p.inventory {
		description += fmt.Sprintf("\n - %s: %s", item.Name, item.InInvDesc)
	}
	return description
}

type quest struct {
	d      *dungeon.Dungeon
	p      *Player
	cellID string

	// allow for us to play over stdout/stdin by default
	// but allow for us to easily change this to work over
	// other protocols, like http, email, or sms
	w io.Writer
	r io.Reader
}

// NewQuest ....
func NewQuest(p *Player, d *dungeon.Dungeon) *quest {
	return &quest{p: p, d: d, cellID: "", w: os.Stdout, r: os.Stdin}
}

func (q *quest) Writer(w io.Writer) {
	q.w = w
}

func (q *quest) Reader(r io.Reader) {
	q.r = r
}

type stateCommand struct {
	currentCellID string
	prompt        string
}

func (q *quest) Start() {
	cmds := make(chan stateCommand, 1)
	go func() {
		// first load, will default to starting cell
		cell, _ := q.d.LoadCell("")
		cmds <- stateCommand{currentCellID: cell.ID, prompt: cell.Prompt("> ")}
	}()
	var exitApp bool

	for !exitApp {
		select {
		case cmd := <-cmds:
			utils.Debugf("got a command for cell " + cmd.currentCellID)
			q.w.Write([]byte(cmd.prompt))
			reader := bufio.NewReader(q.r)
			line, _, err := reader.ReadLine()
			if err != nil {
				utils.Debugf("error reading input - %v", err)
				cmds <- stateCommand{currentCellID: cmd.currentCellID, prompt: "there was an error: " + err.Error()}
				break
			}
			utils.Debugf("received input: " + string(line))
			parsed := parser.Parse(string(line))
			utils.Debugf("%#v", parsed)

			cell, ok := q.d.LoadCell(cmd.currentCellID)
			if !ok {
				utils.Debugf("landed in a bad cell!")
				cmds <- stateCommand{currentCellID: cell.ID, prompt: "you feel a strange sensation, you are suddenly not where you were\n> "}
				break
			}

			if parsed.Action == "look" && (parsed.Object == "" || parsed.Object == "around") {
				cmds <- stateCommand{currentCellID: cell.ID, prompt: cell.Prompt("> ")}
				break
			}

			if parsed.Action == "inventory" || parsed.Action == "look" && parsed.Object == "inventory" {
				cmds <- stateCommand{currentCellID: cell.ID, prompt: q.p.DescribeInventory() + "\n> "}
				break
			}

			// TODO - this should not allow sms to cause the running app to exit
			if parsed.Action == "exit" {
				exitApp = true
				break
			}

			if parsed.Action == "help" {
				cmds <- stateCommand{currentCellID: cell.ID, prompt: helpDialog() + "\n> "}
				break
			}

			describedObject := strings.TrimSpace(parsed.Identifier + " " + parsed.Object)

			if parsed.Action == "go" {
				nextCellID, ok := cell.GetDestinationID(parsed.Object)
				if !ok {
					// maybe "go through the door" did not work, but "go through the green door" will.
					nextCellID, ok = cell.GetDestinationID(describedObject)
					if !ok {
						cmds <- stateCommand{currentCellID: cell.ID, prompt: "hm. That did not work.\n> "}
						break
					}
				}
				// check if the path is blocked by a door

				cellDoor, ok := cell.GetDoor(describedObject)
				if !ok || cellDoor.IsOpen {
					// do door blocking
					cell, _ := q.d.LoadCell(nextCellID)
					cmds <- stateCommand{currentCellID: nextCellID, prompt: cell.Prompt("> ")}
					break
				}
				// door blocks the way
				if cellDoor.IsLocked {
					cmds <- stateCommand{currentCellID: cell.ID, prompt: "The door is locked.\n> "}
					break
				}
				if !cellDoor.IsOpen {
					cmds <- stateCommand{currentCellID: cell.ID, prompt: "The door is not open.\n> "}
					break
				}
			}

			// do door action?
			if cellDoor, ok := cell.GetDoor(describedObject); ok {
				cmds <- stateCommand{currentCellID: cell.ID, prompt: cellDoor.PerformActionAndPrompt(cell, parsed, q.p.inventory...) + "\n> "}
				break
			}

			cmds <- stateCommand{currentCellID: cmd.currentCellID, prompt: "** you can't do that. You must be more specific. Type 'help' to get an idea of how to interact. **\nthis incident has been logged by the system administrator\n> "}
		}
	}
}

func helpDialog() string {
	return "This in an interactive game. You can give simple commands like 'inventory', 'look', 'go north', 'open the yellow chest with the green key' or similar. The most advanced commands are 'verb adjective noun with the adjective noun'."
}
