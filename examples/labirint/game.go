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
	graphic          bool
	logBuf           strings.Builder

	steps       int
	moveHistory []control
}

func newGame(level [][]gameBlock, debug, graphic bool, round int) *game {
	gm := &game{round: round, playerX: -1, playerY: -1, WINX: -1, WINY: -1, debug: debug, graphic: graphic, logBuf: strings.Builder{}}

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
	if gm.graphic {
		gm.clearScreen()

		for _, row := range gm.level {
			for _, blk := range row {
				fmt.Printf("%s", string(blk.symbol))
			}
			fmt.Print("\n")
		}
		fmt.Print("\033[u\033[K")
	}

	log.Printf("[%d/%d] %s", gm.round, gm.steps, gm.logBuf.String())
	gm.logReset()
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
	res := strconv.Itoa(gm.playerX) + ":" + strconv.Itoa(gm.playerY)
	return res
	buf := strings.Builder{}
	//if gm.win {
	//	buf.WriteString(blokWIN.symbol + "~")
	//}
	for k := range gm.moveHistory {
		buf.WriteString(gm.moveHistory[k].String())
		//buf.WriteString("~")
	}
	//buf.WriteString("~" + strconv.FormatFloat(gm.distToWin(gm.playerX, gm.playerY), 'f', 4, 64))
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
	if gm.win {
		return nil
	}
	res := make([]qlearning.Action, 0, 4)
	res = append(res, &Step{
		gm:          gm,
		Dir:         top,
		BlockAtStep: gm.getBlockFromPlayer(top),
		DistToWin:   gm.distToWin(gm.playerX+top.x, gm.playerY+top.y),
	})
	res = append(res, &Step{
		gm:          gm,
		Dir:         bottom,
		BlockAtStep: gm.getBlockFromPlayer(bottom),
		DistToWin:   gm.distToWin(gm.playerX+bottom.x, gm.playerY+bottom.y),
	})
	res = append(res, &Step{
		gm:          gm,
		Dir:         left,
		BlockAtStep: gm.getBlockFromPlayer(left),
		DistToWin:   gm.distToWin(gm.playerX+left.x, gm.playerY+left.y),
	})
	res = append(res, &Step{
		gm:          gm,
		Dir:         right,
		BlockAtStep: gm.getBlockFromPlayer(right),
		DistToWin:   gm.distToWin(gm.playerX+right.x, gm.playerY+right.y),
	})
	return res
}

func (gm *game) Stat() (bool, int) {
	return gm.win, gm.steps
}
func (gm *game) Move(control control) {
	gm.steps++

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
	gm.l("move %s", string(control.dir))
	// check collider
	if !gm.level[newx][newy].canGoThought {
		gm.l("cant move thought, %T", gm.level[newx][newy])
		return
	}

	gm.moveHistory = append(gm.moveHistory, control)
	if gm.level[newx][newy].isWin {
		gm.win = true
	}

	gm.level[newx][newy] = blokPrs
	gm.level[gm.playerX][gm.playerY] = blokEpt

	gm.playerX = newx
	gm.playerY = newy

}
func (gm *game) logReset() {
	gm.logBuf.Reset()
}
func (gm *game) l(str string, args ...interface{}) {
	if gm.debug {
		gm.logBuf.WriteString(fmt.Sprintf(str, args...))
		gm.logBuf.WriteRune(',')
	}
}

type Refere struct {
	baseScore float64
	stepLimit int
}

func NewRefere(stepLimit int) *Refere {
	return &Refere{stepLimit: stepLimit, baseScore: 1}
}

// Reward calculate effectivity of choosed steps
func (r *Refere) Reward(action *qlearning.StateAction) float64 {
	if st, ok := action.Action.(*Step); ok {
		if !st.BlockAtStep.canGoThought {
			return r.baseScore * -1.0
		}
	}

	gm, ok := action.State.(*game)
	if !ok {
		panic("cant cast state to game")
	}
	if gm.win {
		//return r.baseScore * (float64(r.stepLimit - gm.steps))
	}

	if gm.steps >= r.stepLimit {
		//return r.baseScore * -1.0
	}
	dst := gm.distToWin(gm.playerX, gm.playerY)
	if dst == 0 {
		return r.baseScore
	}
	if math.IsInf(dst, 1) {
		fmt.Println("!")
	}
	ret := r.baseScore / float64(dst) //(float64(gm.steps)) //* float64(gm.distToWin(gm.playerX, gm.playerY)))
	if math.IsInf(ret, 1) {
		fmt.Println("!")
	}
	return ret
}
