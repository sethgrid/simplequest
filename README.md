# Simple Text Quest

This is a demo of interactive texting with Twilio's api

To complete the mission, you must have a `SENDGRID_API_KEY` exported as an env var, but you can play without that. When you have to speak the secret word to the stone table, you can cheat and find the magic word in the source code.

# playing the game

text quest responds to simple commands, usually in the form of `verb noun`, `verb adjective noun`, or, in some cases, `verb adjective noun with adjective noun`. For example, `go east`, `open the blue door`, or `unlock the red door with the blue key`.

# running the game
Run the server by going to the `cmd/quest` directory and run `go run main.go`. You can optionally send it the `-debug` flag.

Command line mode: `go run main.go`

HTTP mode: `go run main.go --mode=http`

HTTP mode allows the endpoints `/sms` and `/plain` to be exposed. Configure your twilio sms enabled number to have a webhook of the `/sms` path to have a interactive text adventure that costs money :)