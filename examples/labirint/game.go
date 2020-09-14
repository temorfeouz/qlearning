package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/temorfeouz/qlearning"
)

type game struct {
	round int

	level            [][]gameBlock
	win              bool
	playerX, playerY int
	WINX, WINY       int
	debug            bool
	logBuf           strings.Builder
	moveHistory      []gameBlock
}

func newGame(level [][]gameBlock, debug bool, round int) *game {
	gm := &game{round: round, playerX: -1, playerY: -1, WINX: -1, WINY: -1, debug: debug, logBuf: strings.Builder{}}

	gm.level = make([][]gameBlock, len(level))

	// search player
	for rowInd, row := range level {
		gm.level[rowInd] = make([]gameBlock, len(level[rowInd]))
		for k, block := range row {
			gm.level[rowInd][k] = block
			if block == blokPrs {
				gm.playerX = rowInd
				gm.playerY = k
			}
			if block == blokWIN {
				gm.WINX = rowInd
				gm.WINY = k
			}
		}
	}
	if gm.playerX == -1 {
		panic("cant find player on level")
	}
	if gm.WINX == -1 {
		panic("cant find WIN on level")
	}
	gm.l("new game started")
	return gm
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
	log.Printf("[%d/%d] %s", gm.round, len(gm.moveHistory), gm.logBuf.String())
	gm.logBuf.Reset()
	fmt.Print(strings.Repeat("\033[A", len(gm.level)+1)) // move the cursor up
}
func (gm *game) clearScreen() {
	// clear
	for range gm.level {
		fmt.Println("\033[u\033[K")
	}
	fmt.Print(strings.Repeat("\033[A", len(gm.level)+1)) // move the cursor up
}

func (gm *game) String() string {
	buf := strings.Builder{}
	if gm.win {
		buf.WriteString(blokWIN.symbol + "~")
	}
	for k := range gm.moveHistory {
		buf.WriteString(gm.moveHistory[k].String())
		buf.WriteString("~")
	}
	buf.WriteString("~" + strconv.FormatFloat(gm.distToWin(gm.playerX, gm.playerY), 'f', 4, 64))
	return buf.String()
}

func (gm *game) getBlockFromPlayer(c control) gameBlock {
	return gm.level[gm.playerX+c.x][gm.playerY+c.y]
}

//dist for block
func (gm *game) distToWin(x, y int) float64 {
	return math.Sqrt(float64(gm.WINX*gm.WINX - x*x + gm.WINY*gm.WINY - y*y))
}

// Next collect all moves from current positions
func (gm *game) Next() []qlearning.Action {
	res := make([]qlearning.Action, 0, 4)
	res = append(res, &Step{
		Dir:         top,
		BlockAtStep: gm.getBlockFromPlayer(top),
		DistToWin:   gm.distToWin(gm.playerX+top.x, gm.playerY+top.y),
	})
	res = append(res, &Step{
		Dir:         bottom,
		BlockAtStep: gm.getBlockFromPlayer(bottom),
		DistToWin:   gm.distToWin(gm.playerX+bottom.x, gm.playerY+bottom.y),
	})
	res = append(res, &Step{
		Dir:         left,
		BlockAtStep: gm.getBlockFromPlayer(left),
		DistToWin:   gm.distToWin(gm.playerX+left.x, gm.playerY+left.y),
	})
	res = append(res, &Step{
		Dir:         right,
		BlockAtStep: gm.getBlockFromPlayer(right),
		DistToWin:   gm.distToWin(gm.playerX+right.x, gm.playerY+right.y),
	})
	return res
}

func (gm *game) Stat() (bool, int) {
	return gm.win, len(gm.moveHistory)
}
func (gm *game) Move(control control) {
	gm.moveHistory = append(gm.moveHistory, gm.getBlockFromPlayer(control))

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

// deprecated
type whatISee struct {
	Right, Left, Top, Bottom gameBlock
	DistToWin                float64
}

// deprecated
func (gm *game) LookAround() whatISee {
	res := whatISee{}
	res.Top = gm.level[gm.playerX-1][gm.playerY]
	res.Bottom = gm.level[gm.playerX+1][gm.playerY]
	res.Left = gm.level[gm.playerX][gm.playerY-1]
	res.Right = gm.level[gm.playerX][gm.playerY+1]
	res.DistToWin = math.Sqrt(float64(gm.WINX*gm.WINX - gm.playerX*gm.playerX + gm.WINY*gm.WINY - gm.playerY*gm.playerY))
	gm.l("to win:%+v", res.DistToWin)
	return res
}

type Refere struct {
	baseScore float64
	stepLimit int
	counter   int
}

func NewRefere(stepLimit int) *Refere {
	return &Refere{stepLimit: stepLimit, baseScore: 1000}
}
func (r *Refere) Inc() {
	r.counter++
}
func (r *Refere) Reset() {
	r.counter = 0
}

// Reward calculate effectivity of choosed steps
func (r *Refere) Reward(action *qlearning.StateAction) float32 {
	if strings.Contains(action.State.String(), blokWIN.String()) {
		return float32(r.baseScore) * 10.
	}

	tmp := strings.Split(action.State.String(), "~")

	dist, err := strconv.ParseFloat(tmp[len(tmp)-1], 64)
	if err != nil {
		panic(err)
	}
	if r.counter >= r.stepLimit {
		return float32(r.baseScore * -1.0)
	}

	return float32(r.baseScore*dist) / float32(len(tmp)-1)
}
