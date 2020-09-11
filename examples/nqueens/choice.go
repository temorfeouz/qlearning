package main

import (
	"strconv"
	"strings"

	"github.com/ecooper/qlearning"
)

// Choice implements qlearning.Action for a character choice in a game
// of Hangman.
type Choice struct {
	variants []choiceElem
}
type choiceElem struct {
	row    int
	column int
}

// String returns the character for the current action.
func (choice *Choice) String() string {
	str := strings.Builder{}
	for i := range choice.variants {
		str.WriteString(strconv.Itoa(choice.variants[i].row) + ":" + strconv.Itoa(choice.variants[i].column) + ";")
	}
	return str.String()
}

// Apply updates the state of the game for a given character choice.
func (choice *Choice) Apply(state qlearning.State) qlearning.State {
	game := state.(*Game)
	game.Choose(choice)

	return game
}
