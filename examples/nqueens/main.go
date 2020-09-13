package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/temorfeouz/qlearning"
)

const (
	Lost   = -1
	Active = 0
	Won    = 1

	//NQUEENS = 15
	NQUEENS    = 8
	BOARD_SIZE = 10
	qtableFile = "qtable.json"
)

var (
	progressAt int = 1000
	playFor    int = 5000000
)

func init() {
	rand.Seed(time.Now().UnixNano())

	if !fileExists(qtableFile) {
		f, err := os.Create(qtableFile)
		if err != nil {
			panic(err)
		}
		f.Close()
	}
}
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func init() {
	//flag.StringVar(&wordListPath, "wordlist", wordListPath, "Path to a wordlist")
	//flag.BoolVar(&debug, "debug", debug, "Set debug")
	flag.IntVar(&progressAt, "progress", progressAt, "Print progress messages every N games")
	//flag.IntVar(&wordCount, "words", wordCount, "Use N words from wordlist")
	flag.IntVar(&playFor, "games", playFor, "Play N games")

	flag.Parse()
}

func main() {
	var (
		collisions float32 = 0.0
		count              = 0
		board      string

		// Our agent has a learning rate of 0.7 and discount of 1.0.
		agent = qlearning.NewSimpleAgent(0.7, 1.0)
	)

	f, err := os.OpenFile(qtableFile, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	agent.Import(f)
	f.Close()

	progress := func() {
		// Print our progress every 1000 rows.
		if count > 0 && count%progressAt == 0 {
			fmt.Printf("%d games played, collisions count %f \n board is:\n%s\n", count, collisions, board)
		}
	}

	// Let's play 5 million games
	for count = 1; count < playFor; count++ {
		// Get a new word and game for each iteration...
		game := NewGame(NQUEENS)

		//game.Log("Game created")

		// While the game is still active, we'll continue to update
		// our agent and learn from its choices.
		for game.IsComplete() == 0 {
			// Pick the next move, which is going to be a letter choice.
			action := qlearning.Next(agent, game, 0.1)

			// Whatever that choice is, let's update our model for its
			// impact. If the character chosen is in the game's word,
			// then this action will be positive. Otherwise, it will be
			// negative.
			agent.Learn(action, game)
			board = game.DrawBoard()
			collisions = game.collisions
			if collisions == 0 {
				game.Log("FINISH ON STEP %d", count)
				game.Log(game.DrawBoard())
				log.Println("store progress...")
				store(agent)
				log.Println("progress finished")
				return
			} else {
				//game.Log("[STEP %d] got %f collisions", count, collisions)
				//game.Log(game.String())
			}
		}

		progress()
	}

	progress()
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
func store(agent *qlearning.SimpleAgent) {
	f, err := os.OpenFile(qtableFile, os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	agent.Export(f)
}
