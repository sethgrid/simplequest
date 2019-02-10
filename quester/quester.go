package quester

import (
	"fmt"
	"strings"
	"time"

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

type Quest struct {
	d      *dungeon.Dungeon
	p      *Player
	cellID string

	lastCommand time.Time

	in   chan string
	out  chan string
	cmds chan stateCommand

	isStopped bool
}

// NewQuest ....
func NewQuest(p *Player, d *dungeon.Dungeon) *Quest {
	return &Quest{p: p, d: d, cellID: "", in: make(chan string), out: make(chan string)}
}

type stateCommand struct {
	currentCellID string
	prompt        string
}

func (q *Quest) TakeCommand(s string) string {
	go func() { q.in <- s }()
	return <-q.out
}

var initilizer = "49hhkjndsf94"

func (q *Quest) IsExpired() bool {
	return q.lastCommand.Add(2 * time.Minute).Before(time.Now())
}

func (q *Quest) Stop() {
	q.isStopped = true
	close(q.cmds)
}

func (q *Quest) Start() {
	q.cmds = make(chan stateCommand, 1)

	go func() {
		// on the first game load over http/sms, we get our first input and must discard it
		// because they have not gotten the first game prompt
		cell, _ := q.d.LoadCell("")
		<-q.in
		q.cmds <- stateCommand{currentCellID: cell.ID, prompt: cell.Prompt("> ")}
	}()

	for !q.isStopped {
		select {
		case cmd := <-q.cmds:
			utils.Debugf("got a command for cell " + cmd.currentCellID)
			q.out <- cmd.prompt
			line := <-q.in

			q.lastCommand = time.Now()

			utils.Debugf("received input: " + line)
			parsed := parser.Parse(line)
			utils.Debugf("%#v", parsed)

			describedObject := strings.TrimSpace(parsed.Identifier + " " + parsed.Object)

			cell, ok := q.d.LoadCell(cmd.currentCellID)
			if !ok {
				utils.Debugf("landed in a bad cell!")
				q.cmds <- stateCommand{currentCellID: cell.ID, prompt: "you feel a strange sensation, you are suddenly not where you were\n> "}
				break
			}

			if parsed.Action == "look" && (parsed.Object == "" || parsed.Object == "around") {
				q.cmds <- stateCommand{currentCellID: cell.ID, prompt: cell.Prompt("> ")}
				break
			}

			if parsed.Action == "look" && describedObject != "" {
				item, ok := cell.GetItem(describedObject)
				if !ok {
					q.cmds <- stateCommand{currentCellID: cell.ID, prompt: fmt.Sprintf("there is no %s\n> ", describedObject)}
					break
				}
				q.cmds <- stateCommand{currentCellID: cell.ID, prompt: fmt.Sprintf("%s\n> ", item.InRoomDesc)}
				break
			}

			if parsed.Action == "inventory" || parsed.Action == "look" && parsed.Object == "inventory" {
				q.cmds <- stateCommand{currentCellID: cell.ID, prompt: q.p.DescribeInventory() + "\n> "}
				break
			}

			// TODO - this should not allow sms to cause the running app to exit
			if parsed.Action == "exit" {
				q.Stop()
				break
			}

			if parsed.Action == "help" {
				q.cmds <- stateCommand{currentCellID: cell.ID, prompt: helpDialog() + "\n> "}
				break
			}

			if parsed.Action == "take" {
				item, ok := cell.GetItem(describedObject)
				if !ok {
					q.cmds <- stateCommand{currentCellID: cell.ID, prompt: fmt.Sprintf("there is no %s to take", describedObject) + "\n> "}
					break
				}
				if !item.Takable {
					q.cmds <- stateCommand{currentCellID: cell.ID, prompt: fmt.Sprintf("you cannot take the %s", describedObject) + "\n> "}
					break
				}
				item.InInventory = true
				cell.RemoveItem(item.Name)
				q.p.AddInventory(item)
				q.cmds <- stateCommand{currentCellID: cell.ID, prompt: fmt.Sprintf("you've taken the %s", describedObject) + "\n> "}
				break
			}

			if parsed.Action == "go" {
				nextCellID, ok := cell.GetDestinationID(parsed.Object)
				if !ok {
					// maybe "go through the door" did not work, but "go through the green door" will.
					nextCellID, ok = cell.GetDestinationID(describedObject)
					if !ok {
						q.cmds <- stateCommand{currentCellID: cell.ID, prompt: "hm. That did not work.\n> "}
						break
					}
				}
				// check if the path is blocked by a door

				cellDoor, ok := cell.GetDoor(describedObject)
				if !ok || cellDoor.IsOpen {
					// do door blocking
					cell, _ := q.d.LoadCell(nextCellID)
					q.cmds <- stateCommand{currentCellID: nextCellID, prompt: cell.Prompt("> ")}
					break
				}
				// door blocks the way
				if cellDoor.IsLocked {
					q.cmds <- stateCommand{currentCellID: cell.ID, prompt: "The door is locked.\n> "}
					break
				}
				if !cellDoor.IsOpen {
					q.cmds <- stateCommand{currentCellID: cell.ID, prompt: "The door is not open.\n> "}
					break
				}
			}

			// do door action?
			if cellDoor, ok := cell.GetDoor(describedObject); ok {
				q.cmds <- stateCommand{currentCellID: cell.ID, prompt: cellDoor.PerformActionAndPrompt(cell, parsed, q.p.inventory...) + "\n> "}
				break
			}

			// attempt the action on all items in the room
			var itemResponse string
			for _, item := range cell.Items() {
				utils.Debugf("inspecting " + item.Name)
				itemResponse = item.Action(cell, parsed, q.p.inventory...)
				if itemResponse == "nothing happens" || itemResponse == "" {
					utils.Debugf("item does nothing")
					continue
				}
				utils.Debugf("got item response")
				break
			}

			if itemResponse != "" {
				q.cmds <- stateCommand{currentCellID: cell.ID, prompt: itemResponse + "\n> "}
				break
			}

			q.cmds <- stateCommand{currentCellID: cmd.currentCellID, prompt: "** you can't do that. You must be more specific. Type 'help' to get an idea of how to interact. **\nthis incident has been logged by the system administrator\n> "}
		}
	}
}

func helpDialog() string {
	return "This in an interactive game. You can give simple commands like 'inventory', 'look', 'go north', 'open the yellow chest with the green key' or similar. The most advanced commands are 'verb adjective noun with the adjective noun'."
}
