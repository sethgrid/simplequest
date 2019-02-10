package quester

import (
	"fmt"
	"log"
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

func (q *Quest) IsExpired() bool {
	return q.lastCommand.Add(2 * time.Hour).Before(time.Now())
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
		q.cmds <- stateCommand{currentCellID: cell.ID, prompt: cell.Prompt("")}
	}()

	for !q.isStopped {
		select {
		case cmd := <-q.cmds:
			q.lastCommand = time.Now()
			q.out <- cmd.prompt
			line := <-q.in
			parsed := parser.Parse(line)
			cell, ok := q.d.LoadCell(cmd.currentCellID)
			if !ok {
				log.Printf("bad path exists: current cell %#v", cmd.currentCellID)
				q.cmds <- stateCommand{currentCellID: cell.ID, prompt: "you feel a strange sensation, you are suddenly not where you were "}
				break
			}

			describedObject := strings.TrimSpace(parsed.Identifier + " " + parsed.Object)
			utils.Debugf("%q [%s] %#v - described object %q", q.p.id, cmd.currentCellID, parsed, describedObject)

			if parsed.Action == "look" && (parsed.Object == "" || parsed.Object == "around") {
				q.cmds <- stateCommand{currentCellID: cell.ID, prompt: cell.Prompt("")}
				break
			}

			if parsed.Action == "look" && describedObject != "" {
				item, ok := cell.GetItem(describedObject)
				if !ok {
					q.cmds <- stateCommand{currentCellID: cell.ID, prompt: fmt.Sprintf("there is no %s ", describedObject)}
					break
				}
				q.cmds <- stateCommand{currentCellID: cell.ID, prompt: fmt.Sprintf("%s ", item.InRoomDesc)}
				break
			}

			if parsed.Action == "inventory" || parsed.Action == "look" && parsed.Object == "inventory" {
				q.cmds <- stateCommand{currentCellID: cell.ID, prompt: q.p.DescribeInventory()}
				break
			}

			// TODO - this should not allow sms to cause the running app to exit
			if parsed.Action == "exit" {
				q.Stop()
				q.out <- "you have existed text quest"
				break
			}

			if parsed.Action == "help" || parsed.Action == "manual" {
				q.cmds <- stateCommand{currentCellID: cell.ID, prompt: helpDialog()}
				break
			}

			if parsed.Action == "take" {
				item, ok := cell.GetItem(describedObject)
				if !ok {
					q.cmds <- stateCommand{currentCellID: cell.ID, prompt: fmt.Sprintf("there is no %s to take", describedObject)}
					break
				}
				if !item.Takable {
					q.cmds <- stateCommand{currentCellID: cell.ID, prompt: fmt.Sprintf("you cannot take the %s", describedObject)}
					break
				}
				item.InInventory = true
				cell.RemoveItem(item.Name)
				q.p.AddInventory(item)
				q.cmds <- stateCommand{currentCellID: cell.ID, prompt: fmt.Sprintf("you've taken the %s", describedObject)}
				break
			}

			if parsed.Action == "go" {
				nextCellID, ok := cell.GetDestinationID(parsed.Object)
				if !ok {
					// maybe "go through the door" did not work, but "go through the green door" will.
					utils.Debugf("go action. Next cell id is %q (parsed object %q)", nextCellID, parsed.Object)
					nextCellID, ok = cell.GetDestinationID(describedObject)
					utils.Debugf("go action. Next cell id is %q (parsed descirbed object %q)", nextCellID, describedObject)
					if !ok {
						utils.Debugf("go action wont work")
						q.cmds <- stateCommand{currentCellID: cell.ID, prompt: "hm. That did not work. "}
						break
					}
				}
				// check if the path is blocked by a door
				utils.Debugf("check door actions for %q", describedObject)
				cellDoor, ok := cell.GetDoor(describedObject)
				if !ok || cellDoor.IsOpen {
					// do door blocking
					cell, _ := q.d.LoadCell(nextCellID)
					q.cmds <- stateCommand{currentCellID: nextCellID, prompt: cell.Prompt("")}
					break
				}
				// door blocks the way
				if cellDoor.IsLocked {
					q.cmds <- stateCommand{currentCellID: cell.ID, prompt: "The door is locked. Perhaps it needs to be unlocked. "}
					break
				}
				if !cellDoor.IsOpen {
					q.cmds <- stateCommand{currentCellID: cell.ID, prompt: fmt.Sprintf("The door is not open. Try 'open the %s ", cellDoor.Name)}
					break
				}
			}

			// do door action?
			if cellDoor, ok := cell.GetDoor(describedObject); ok {
				q.cmds <- stateCommand{currentCellID: cell.ID, prompt: cellDoor.PerformActionAndPrompt(cell, parsed, q.p.inventory...)}
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
				q.cmds <- stateCommand{currentCellID: cell.ID, prompt: itemResponse}
				break
			}

			utils.Debugf("unrecognized command:\n  player %#v\n  cell %#v\n  parsed  %#v", q.p, cell, parsed)
			q.cmds <- stateCommand{currentCellID: cmd.currentCellID, prompt: "Unrecognized command. This incident has been logged by the system administrator.\n"}
		}
	}
}

func helpDialog() string {
	return "This in an interactive game. You can give simple commands like 'inventory', 'look', 'go north', 'open the yellow chest with the green key' or similar. The most advanced commands are 'verb adjective noun with the adjective noun'."
}
