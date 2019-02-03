package main

import (
	"github.com/sethgrid/simplequest/dungeon"
)

func main() {
	// get config, run game
	d := dungeon.MakeDungeon()
	c := d.NewCell("0,1")
	c.Description("you are in a clearing")

}
