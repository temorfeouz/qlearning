package main

import (
	"fmt"
	"log"
	"strings"
)

type game struct {
	level   [][]gameBlock
	win     bool
	steps   int64
	playerX int
	playerY int
	debug   bool
	logBuf  strings.Builder
}

func newGame(level [][]gameBlock, debug bool) *game {
	gm := &game{level: level, playerX: -1, playerY: -1, debug: debug, logBuf: strings.Builder{}}
	// search player
	for rowInd, row := range level {
		for k, block := range row {
			if block == blokPrs {
				gm.playerX = rowInd
				gm.playerY = k
				break
			}
		}
	}
	if gm.playerX == -1 {
		panic("cant find player on level")
	}
	gm.l("new game started")
	return gm
}
func (gm *game) clearScreen() {
	// clear
	for range gm.level {
		fmt.Println("\033[u\033[K")
	}
	fmt.Print(strings.Repeat("\033[A", len(gm.level)+1)) // move the cursor up

}
func (gm *game) Draw() {
	gm.clearScreen()

	for _, row := range gm.level {
		for _, blk := range row {
			fmt.Printf("%s", string(blk.symbol))
		}
		fmt.Print("\n")
	}

	fmt.Print("\033[u\033[K")
	log.Print(gm.logBuf.String())
	gm.logBuf.Reset()
	fmt.Print(strings.Repeat("\033[A", len(gm.level)+1)) // move the cursor up
	//for {w
	//
	//	fmt.Print(strings.Repeat("\033[A", len(gm.level))) // move the cursor up
	//	//fmt.Printf("Retrieved %d\n", i)
	//	time.Sleep(time.Second)
	//
	//}
}
func (gm *game) Stat() (bool, int64) {
	return gm.win, gm.steps
}
func (gm *game) Move(control control) {
	gm.steps += 1

	newx := gm.playerX + control.x
	newy := gm.playerY + control.y

	// isCanMove
	if newx >= len(gm.level) || newx < 0 {
		gm.l("cant move in that direction(x)")
		return
	} else if newy >= len(gm.level[newx]) || newy < 0 {
		gm.l("cant move in that direction(y)")
		return
	}

	// check collider
	if !gm.level[newx][newy].canGoThought {
		gm.l("cant move thought, %T", gm.level[newx][newy])
		return
	} else if gm.level[newx][newy].isWin {
		gm.win = true
	}

	gm.level[newx][newy] = blokPrs
	gm.level[gm.playerX][gm.playerY] = blokEpt

	gm.playerX = newx
	gm.playerY = newy

}
func (gm *game) l(str string, args ...interface{}) {
	if gm.debug {
		gm.logBuf.WriteString(fmt.Sprintf(str, args...))
		gm.logBuf.WriteRune(',')
	}
}
