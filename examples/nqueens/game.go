package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/temorfeouz/qlearning"
)

// Game represents the state of any given game of Hangman. It implements
// qlearning.Agent, qlearning.Rewarder, and qlearning.State.
type Game struct {
	queenCount int
	desk       []int
	collisions float32

	debug bool
}

// NewGame creates a new Hangman game for the given word. If debug
// is true, Game.Log messages will print to stdout.
func NewGame(queensCount int) *Game {
	game := &Game{debug: true}
	game.New(queensCount)

	return game
}

// New resets the current game to a new game for the given word.
func (game *Game) New(queenCount int) {
	game.queenCount = queenCount
	//game.Lives = queenCount
	game.desk = make([]int, queenCount)

	rand.Seed(time.Now().UnixNano())
	for i := range game.desk {
		game.desk[i] = rand.Intn(game.queenCount - 1)
	}
}

// Returns Lost, Active, or Won based on the game's current state.
func (game *Game) IsComplete() int {
	if game.collisions == float32(game.queenCount) {
		return Won
	}

	return Active
}

// Choose applies a character attempt in the current game, returning
// true if char is present in Game.Word.
//
// Choose updates the game's state.
func (game *Game) Choose(choice *Choice) bool {
	for _, c := range choice.variants {
		game.desk[c.row] = c.column
	}
	return true
}

// Reward returns a score for a given qlearning.StateAction. Reward is a
// member of the qlearning.Rewarder interface.
func (game *Game) Reward(action *qlearning.StateAction) float32 {
	//ss := action.State.(*Game)
	game.collisions = 0
	for i := 0; i < len(game.desk); i++ {
		// search right angle collisions
		for j := i + 1; j < len(game.desk); j++ {
			if j-i == absInt(game.desk[i]-game.desk[j]) {
				game.collisions++
			}
		}
	}
	return float32(game.queenCount) - game.collisions
}

// Next creates a new slice of qlearning.Action instances. A possible
// action is created for each character that has not been attempted in
// in the game.
func (game *Game) Next() []qlearning.Action {
	actions := make([]qlearning.Action, 1)

	allowedVals := make([]int, game.queenCount)
	tmp := &Choice{}
	for i := 0; i < game.queenCount; i++ {
		allowedVals[i] = i
		tmp.variants = append(tmp.variants, choiceElem{row: i})
	}
	rand.Shuffle(len(allowedVals), func(i, j int) { allowedVals[i], allowedVals[j] = allowedVals[j], allowedVals[i] })
	for k, v := range allowedVals {
		tmp.variants[k].column = v
	}

	actions[0] = tmp

	return actions
}

// Log is a wrapper of fmt.Printf. If Game.debug is true, Log will print
// to stdout.
func (game *Game) Log(msg string, args ...interface{}) {
	if game.debug {
		logMsg := fmt.Sprintf("%s\n", msg)
		fmt.Printf(logMsg, args...)
	}
}

// String returns a consistent hash for the current game state to be
// used in a qlearning.Agent.
func (game *Game) String() string {
	var board bytes.Buffer
	for k, v := range game.desk {
		board.WriteString(strconv.Itoa(k) + ":" + strconv.Itoa(v) + ";")
	}
	return board.String()
}
func (game *Game) DrawBoard() string {
	var board bytes.Buffer
	for _, p := range game.desk {
		if p < 0 {
			fmt.Printf("%+v", "!!")
			os.Exit(1)
		}
		board.WriteString(strings.Repeat(" .", p))
		board.WriteString(" \u2655")
		board.WriteString(strings.Repeat(" .", NQUEENS-p-1))
		board.WriteString("\n")
	}
	return board.String()
}
