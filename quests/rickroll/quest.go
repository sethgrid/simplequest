package rickroll

/* Planned Map and room guide */

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

import (
	"fmt"
	"log"
	"os"
	"strings"

	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sethgrid/simplequest/dungeon"
	"github.com/sethgrid/simplequest/parser"
	"github.com/sethgrid/simplequest/quester"
	"github.com/sethgrid/simplequest/utils"
)

var sgAPIKey string

// NewRickRoll sets up the dungeon and loads it into a quest
func NewRickRoll(p *quester.Player) *quester.Quest {
	sgAPIKey, _ = os.LookupEnv("SENDGRID_API_KEY")
	if sgAPIKey == "" {
		log.Fatal("unable to load NewRickRoll quest withough SENDGRID_API_KEY environment variable set")
	}

	d := dungeon.MakeDungeon()
	// opening clearing
	d.NewCell("outworld 1,1").
		Description("you are in a clearing").
		AddDestination("outworld 1,0", "north", "north, a dark forest").
		AddDestination("outworld 2,1", "east", "to the east, large spires peak over the hills")

	// dark forest north of clearing
	d.NewCell("outworld 1,0").
		Description("you find yourself in a dark forest").
		AddDestination("outworld 1,1", "south", "to the south, a clearing")

	// castle courtyard
	d.NewCell("outworld 2,1").
		Description("you stand before a great castle. To the north is dense forest, and to the south, a great chasm").
		AddDestination("outworld 1.1", "west", "to the west lies a clearing").
		AddDestination("castle 0,1", "steps|door|rusty door|east", "Up the large steps stands a large door")

	// castle enterance
	d.NewCell("castle 0,1").
		Description("you arrive at the top of the steps, a little out of breath.").
		AddDestination("castle 1,1", "east|door|rusty door", "a giant rusty door, it is open just enough to fit through.").
		AddDestination("outworld 2,1", "steps|west", "going down the steps to the west leads away from the castle.").
		AddDoor(rustyDoor)

	// castle great entry way
	d.NewCell("castle 1,1").
		Description("you enter the geat entry way. Soon, there will be three doors.").
		AddDestination("castle 0,1", "rusty door", "the rusty door is behind you").
		AddDestination("castle 2,0", "blue door", "a blue door").
		AddDestination("castle 2,1", "red door", "a red door with a blue lock on it").
		AddDestination("castle 2,2", "green door", "a green door with a white lock on it").
		AddDoor(rustyDoor).
		AddDoor(redDoor).
		AddDoor(blueDoor).
		AddDoor(greenDoor)

	// blue door
	d.NewCell("castle 2,0").
		Description("in the center of room is a large stone table").
		AddDestination("castle 1,1", "door|blue door", "the door out is on the eastern wall").
		AddDoor(redDoor).
		AddItem(stoneTable)

	// red door
	d.NewCell("castle 2,1").
		Description("There is a white door. It is sealed shut with no apparent lock, inscribed on its doors is https://bit.ly/IqT6zt").
		AddDestination("castle 1,1", "door|red door", "the door out is on the eastern wall").
		AddDoor(redDoor)

	// green door
	d.NewCell("castle 2,2").
		Description("How did you get in this room, there was no white key...").
		AddDestination("castle 1,1", "door|green door", "the door out is on the eastern wall").
		AddDoor(greenDoor)

	// white room
	d.NewCell("castle 3,2").
		Description("How did you get through to this room, nothing connects to it!").
		AddDestination("castle 1,1", "door|white door", "the door out is on the eastern wall")

	masterkey := &dungeon.Item{
		Name:        "master key",
		Takable:     true,
		InInventory: true,
		InInvDesc:   "it can unlock any door",
	}
	_ = masterkey

	player := quester.NewPlayer("foo")

	quest := quester.NewQuest(player, d)

	return quest
}

var stoneTable = &dungeon.Item{
	Name:       "stone table",
	Takable:    false,
	Movable:    false,
	InRoomDesc: "a large stone table. An incription reads 'tell me first your email, then a secret'",
	Action: func(cell *dungeon.Cell, command parser.Parsed, inventory ...*dungeon.Item) string {
		utils.Debugf("inside the stone table action router")
		var email string
		if strings.Contains(command.RawAction, "@") {
			email = command.RawAction
		} else if strings.Contains(command.Object, "@") {
			email = command.Object
		}
		if email != "" {
			ok := sendMail(email)
			if ok {
				return "the room rumbles for a moment, dust falls from the ceiling. You hear a distant whiper: 'email sent...'"
			}
			return "the begins to rumble, but quickly goes still. You hear distant whisper: 'unable to send email...'"
		}

		if command.Action == "say" && command.Object == "mellon" {
			cell.AddItem(blueKey)
			item, ok := cell.GetItem("stone table")
			if !ok {
				return "something unexpected happened to the stone table. A blue key appeared upon it"
			}
			item.InRoomDesc = "a large stone table lay split in half"
			return "the room violently rumbles and begins to quake. The large stone table strains with the shifting ground and breaks in half. Floating in the air, you see a blue key."
		}

		if command.Object != "table" {
			return "nothing happens"
		}

		if command.Action == "move" {
			return "you are unable to budge the large stone table"
		}

		return "nothing happens"
	},
}

var blueKey = &dungeon.Item{
	Name:       "blue key",
	Takable:    true,
	InRoomDesc: "a blue key floats in midair",
	InInvDesc:  "an ordinary key of the blue persuation - aside from the fact it was found floating in the air.",
	Action: func(cell *dungeon.Cell, command parser.Parsed, inventory ...*dungeon.Item) string {
		switch command.Action {
		case "use":
			return "Use it how? Try: unlock the <thing> with the blue key"
		default:
			return "nothing happens"
		}
	},
}

var rustyDoor = dungeon.NewDoor(
	"rusty door", !dungeon.Locked, dungeon.Open, "the door is very tall, at least 15 ft. The henges are rusted in place, keeping the doors ajar.",
	func(cell *dungeon.Cell, command parser.Parsed, inventory ...*dungeon.Item) string {
		switch command.Action {
		case "push":
			return "push as you might, the door is firmly in place"
		}
		return "nothing happens"
	},
)

// TODO: default lock/unlock and close/open handling with custom descriptions and requried unlocking item
var blueDoor = dungeon.NewDoor(
	"blue door", !dungeon.Locked, !dungeon.Open, "the blue door is firmly shut.",
	func(cell *dungeon.Cell, command parser.Parsed, inventory ...*dungeon.Item) string {
		doorName := strings.TrimSpace(command.Identifier + " " + command.Object)

		switch command.Action {
		case "push", "open":
			if door, ok := cell.GetDoor(doorName); ok {
				if !door.IsOpen {
					door.IsOpen = true
					door.Description = "the red door lays open"
					return "the door easily opens"
				}
				return "the door is already open"
			}
			return fmt.Sprintf("the %s does not seem to be here...", doorName)

		case "close", "shut":
			if door, ok := cell.GetDoor(doorName); ok {
				door.IsOpen = false
				door.Description = "the blue door is closed"

				return "the blue door is closed."
			}
			return fmt.Sprintf("the %s does not seem to be here...", doorName)
		}
		return "nothing happens"
	},
)

var greenDoor = dungeon.NewDoor(
	"green door", dungeon.Locked, !dungeon.Open, "the red door sits firmly in place and is closed. In its center sits a white lock.",
	func(cell *dungeon.Cell, command parser.Parsed, inventory ...*dungeon.Item) string {
		actionObjectName := strings.TrimSpace(command.ActionIdentifier + " " + command.ActionObject)
		doorName := strings.TrimSpace(command.Identifier + " " + command.Object)

		hasActionItem := false
		for _, item := range inventory {
			if item.Name == actionObjectName {
				hasActionItem = true
			}
		}

		switch command.Action {
		case "push", "open":
			if door, ok := cell.GetDoor(doorName); ok {
				if door.IsLocked && !door.IsOpen {
					return "the door does not move, in its center sits a white lock"
				}
				door.IsOpen = true
				door.Description = "the white door lays open"
				return "the door is now open"
			} else {
				return fmt.Sprintf("the %s does not seem to be here...", doorName)
			}
		case "close", "shut":
			if door, ok := cell.GetDoor(doorName); ok {
				door.IsOpen = false
				door.Description = "the red door is closed"
				if door.IsLocked {
					door.Description = "the green door sits firmly in place and is closed. In its center sits a white lock."
				}
				return "the door is no longer open"
			}
			return fmt.Sprintf("the %s does not seem to be here...", doorName)

		case "use", "unlock":
			if actionObjectName == "white key" || actionObjectName == "master key" {
				if !hasActionItem {
					return fmt.Sprintf("you do not have the %s", actionObjectName)
				}
				if door, ok := cell.GetDoor(doorName); ok {
					door.IsLocked = !door.IsLocked
					return "you turn the key and hear a large thunk"
				} else {
					return fmt.Sprintf("the %s does not seem to be here", doorName)
				}
			}
			if actionObjectName == "" {
				return "please be more specific. Unlock the white door with what?"
			}
			return fmt.Sprintf("the %s does nothing, and the %s remains unchanged", actionObjectName, doorName)
		}
		return "nothing happens"
	},
)

var redDoor = dungeon.NewDoor(
	"red door", dungeon.Locked, !dungeon.Open, "the red door sits firmly in place and is closed. In its center sits a blue lock.",
	func(cell *dungeon.Cell, command parser.Parsed, inventory ...*dungeon.Item) string {
		actionObjectName := strings.TrimSpace(command.ActionIdentifier + " " + command.ActionObject)
		doorName := strings.TrimSpace(command.Identifier + " " + command.Object)

		hasActionItem := false
		for _, item := range inventory {
			if item.Name == actionObjectName {
				hasActionItem = true
			}
		}

		switch command.Action {
		case "push", "open":
			if door, ok := cell.GetDoor(doorName); ok {
				if door.IsLocked && !door.IsOpen {
					return "the door does not move, in its center sits a blue lock"
				}
				door.IsOpen = true
				door.Description = "the red door lays open"
				return "the door is now open"
			} else {
				return fmt.Sprintf("the %s does not seem to be here...", doorName)
			}
		case "close", "shut":
			if door, ok := cell.GetDoor(doorName); ok {
				door.IsOpen = false
				door.Description = "the red door is closed"
				if door.IsLocked {
					door.Description = "the red door sits firmly in place and is closed. In its center sits a blue lock."
				}
				return "the door is no longer open"
			}
			return fmt.Sprintf("the %s does not seem to be here...", doorName)

		case "use", "unlock":
			if actionObjectName == "blue key" || actionObjectName == "master key" {
				if !hasActionItem {
					return fmt.Sprintf("you do not have the %s", actionObjectName)
				}
				if door, ok := cell.GetDoor(doorName); ok {
					door.IsLocked = !door.IsLocked
					return "you turn the key and hear a large thunk"
				} else {
					return fmt.Sprintf("the %s does not seem to be here", doorName)
				}
			}
			if actionObjectName == "" {
				return "please be more specific. Unlock the red door with what?"
			}
			return fmt.Sprintf("the %s does nothing, and the %s remains unchanged", actionObjectName, doorName)
		}
		return "nothing happens"
	},
)

func sendMail(emailAddr string) bool {
	from := mail.NewEmail("Text Quest", "textquest@sethammons.com")
	subject := "Message from the stone table"
	to := mail.NewEmail("Text Quest Player", emailAddr)
	plainTextContent := `You've engaged with the stone table in the realm of SMS.
No longer are you confined to the world of text as your world has expanded to email.
Look to the email headers friend, and say the discovered password in the room of the stone table to receive your reward.
`
	htmlContent := strings.Replace(plainTextContent, ".", ".</br>", -1)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	message.Headers = make(map[string]string)
	message.Headers["x-textquest-stone-table-password"] = "mellon"
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
