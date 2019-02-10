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

	// if you define objects outside the scope of this method, they will persist between players :)
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

			if command.Action == "mellon" {
				return "you need a verb. Try 'say mellon'"
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

	var rustyDoor = dungeon.NewDoor(
		"rusty door", !dungeon.Locked, dungeon.Open, "A rusty door stands at least 15 ft tall, the hinges are rusted in place. The rusty door is ajar, enough to pass through in and out of the castle.",
		func(cell *dungeon.Cell, command parser.Parsed, inventory ...*dungeon.Item) string {
			switch command.Action {
			case "push":
				return "push as you might, the door is firmly in place"
			}
			return "nothing happens"
		},
	)

	var whiteDoor = dungeon.NewDoor(
		"white door", dungeon.Locked, !dungeon.Open, "A white door, completely flush with its surroundings and no apparent way to unlock no open it. Inscribed in its center is https://bit.ly/IqT6zt",
		func(cell *dungeon.Cell, command parser.Parsed, inventory ...*dungeon.Item) string {
			return "nothing you do to the white door has any effect."
		},
	)

	var blueDoor = dungeon.NewDoor(
		"blue door", !dungeon.Locked, !dungeon.Open, "A large blue door with no apparent lock. It sits closed. Maybe you can 'open the blue door'?",
		func(cell *dungeon.Cell, command parser.Parsed, inventory ...*dungeon.Item) string {
			doorName := strings.TrimSpace(command.Identifier + " " + command.Object)

			switch command.Action {
			case "push", "open", "kick", "hit", "knock":
				if door, ok := cell.GetDoor(doorName); ok {
					if !door.IsOpen {
						door.IsOpen = true
						door.Description = "A large blue door lays open"
						return "the blue door easily opens"
					}
					return "the blue door is already open"
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
		"green door", dungeon.Locked, !dungeon.Open, "A green door with no hinges is firmly closed. In its center is a large white lock.",
		func(cell *dungeon.Cell, command parser.Parsed, inventory ...*dungeon.Item) string {
			doorName := strings.TrimSpace(command.Identifier + " " + command.Object)

			switch command.Action {
			case "push", "open":
				if door, ok := cell.GetDoor(doorName); ok {
					if door.IsLocked && !door.IsOpen {
						return "the door does not move."
					}
					door.IsOpen = true
					door.Description = "the green door lays open"
					return "the door is now open"
				} else {
					return fmt.Sprintf("the %s does not seem to be here...", doorName)
				}
			case "use", "unlock":
				return fmt.Sprintf("You see no way to unlock the green door.")
			}
			return "nothing happens"
		},
	)

	var redDoor = dungeon.NewDoor(
		"red door", dungeon.Locked, !dungeon.Open, "A massive red door sits firmly in place and is closed. In its center sits a blue lock.",
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
					door.Description = "A red door lays open"
					return "the door is now open"
				} else {
					return fmt.Sprintf("the %s does not seem to be here...", doorName)
				}
			case "close", "shut":
				if door, ok := cell.GetDoor(doorName); ok {
					door.IsOpen = false
					door.Description = "the red door is closed"
					if door.IsLocked {
						door.Description = "A massive red door sits firmly in place and is closed. In its center sits a blue lock."
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
						return "you turn the key and hear a large thunk as the lock slides open."
					} else {
						return fmt.Sprintf("the %s does not seem to be here", doorName)
					}
				}
				if actionObjectName == "" {
					return "please be more specific. Unlock the red door with what?"
				}
				return fmt.Sprintf("the %s does nothing, and the %s remains unchanged", actionObjectName, doorName)

			case "lock":
				if actionObjectName == "blue key" || actionObjectName == "master key" {
					if !hasActionItem {
						return fmt.Sprintf("you do not have the %s", actionObjectName)
					}
					if door, ok := cell.GetDoor(doorName); ok {
						door.IsLocked = true
						return "you turn the key and hear a large thunk as the lock slides into place."
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

	d := dungeon.MakeDungeon()
	// opening clearing
	d.NewCell("outworld 1,1").
		Description("You find yourself in a clearing. You see a sign that reads: \"Welcome to Text Quest! Respond with simple commands and explore this text based world. Type 'manual' if you need more help.\". Maybe you should 'go east'. If you get lost, maybe 'look around'.").
		AddDestination("outworld 1,0", "north", "To the north you see a vast, dark forest.").
		AddDestination("outworld 2,1", "east", "To the east, you see the tops of large spires that just peak over the hills.")

	// dark forest north of clearing
	d.NewCell("outworld 1,0").
		Description("You find yourself in a dark forest. The trees and undergrowth are much too thick to travel through.").
		AddDestination("outworld 1,1", "south", "to the south, you can see a clearing.")

	// castle courtyard
	d.NewCell("outworld 2,1").
		Description("You stand before a great castle. To the north is dense forest (too dense to even approach), and to the south, a great chasm.").
		AddDestination("outworld 1.1", "west", "To the west lies a clearing.").
		AddDestination("castle 0,1", "steps|door|rusty door|east", "Up the large steps stands a large door")

	// castle enterance
	d.NewCell("castle 0,1").
		Description("You stand atop a many, many large steps. This is the entry way into the castle.").
		AddDoor("castle 1,1", "rusty door|east|door|", rustyDoor).
		AddDestination("outworld 2,1", "steps|west", "Going down the steps to the west leads away from the castle.")

	// castle great entry way
	d.NewCell("castle 1,1").
		Description("You are in a great entry way. The ceilings are at least 10 spans above your head. The smell of damp rot tugs at your nose. Three doors lead further into the castle.").
		AddDoor("castle 0,1", "rusty door|east", rustyDoor).
		AddDoor("castle 2,0", "blue door", blueDoor).
		AddDoor("castle 2,1", "red door", redDoor).
		AddDoor("castle 2,2", "green door", greenDoor)

	// blue door
	d.NewCell("castle 2,0").
		Description("The room is quite large, and in the center of room is a large stone table").
		AddDoor("castle 1,1", "door|blue door", blueDoor).
		AddItem(stoneTable)

	// red door
	d.NewCell("castle 2,1").
		Description("The room is a narrow hallway.").
		AddDoor("castle 1,1", "door|red door", redDoor).
		AddDoor("castle 3,2", "white door", whiteDoor)

	// green door
	d.NewCell("castle 2,2").
		Description("How did you get in this room, there was no white key...").
		AddDoor("castle 1,1", "door|green door", greenDoor)

	// white room
	d.NewCell("castle 3,2").
		Description("How did you get through to this room, nothing connects to it!").
		AddDoor("castle 1,1", "door|white door", whiteDoor)

	masterkey := &dungeon.Item{
		Name:        "master key",
		Takable:     true,
		InInventory: true,
		InInvDesc:   "it can unlock any door",
	}
	_ = masterkey

	quest := quester.NewQuest(p, d)

	return quest
}

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
