package main

/*
TODOs
leave room or use door or take key (if there is only one of a thing, they want to use that one of a thing)
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
	flag.StringVar(&gameMode, "mode", "cmd", "set to cmd for command line, http for http server")
	flag.IntVar(&port, "port", 5000, "set the port to run the http server")
	flag.Parse()

	utils.Debug = debug
	if !mode.ValidGameMode(gameMode) {
		log.Fatal("please select a -mode of cmd, http, or sms")
	}

	if gameMode == "http" {
		mode.RunHTTPServer(port)
	}

	mode.RunCMDServer()

}
