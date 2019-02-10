package main

/*
TODO: push open the door. push the door open. go through the <closed> door (give help, you need to open it first)
look for types like geat hall.
start with basic instructions
kick down the door (strip down?)
send an actual email
integrate with twilio sms for convo
change leave (don't exit the app)
leave room or use door or take key (if there is only one of a thing, they want to use that one of a thing)
room descriptions: dont act like you just go in there for look around
door descriptions should say if they are open or closed
lock description - large thunk -> you've unlocked the door
clean up old quests
*/

import (
	"flag"
	"log"

	"github.com/sethgrid/simplequest/mode"
	"github.com/sethgrid/simplequest/utils"
)

func main() {
	// get config, run game
	var debug bool
	var gameMode string
	var port int
	flag.BoolVar(&debug, "debug", false, "set flag to enable debug logs")
	flag.StringVar(&gameMode, "mode", "cmd", "set to cmd for command line, http for http server, or sms for game over text")
	flag.IntVar(&port, "port", 5000, "set the port to run the http server")
	flag.Parse()

	utils.Debug = debug
	if !mode.ValidGameMode(gameMode) {
		log.Fatal("please select a -mode of cmd, http, or sms")
	}

	mode.RunHTTPServer(port)

}
