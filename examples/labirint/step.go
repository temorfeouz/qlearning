package main

import (
	"fmt"
	"strconv"

	"github.com/temorfeouz/qlearning"
)

type Step struct {
	gm          *game
	Dir         control
	BlockAtStep gameBlock
	DistToWin   float64
}

func (s *Step) String() string {
	//res := strings.Builder{}
	//for _, v := range s.gm.moveHistory {
	//	res.WriteString(v.String())
	//}
	//return res.String()
	res := strconv.Itoa(s.gm.playerX+s.Dir.x) + ":" + strconv.Itoa(s.gm.playerY+s.Dir.y)
	return res
	return fmt.Sprintf("cs:%s", s.BlockAtStep.String())
}
func (s *Step) Apply(st qlearning.State) qlearning.State {
	gm, ok := st.(*game)
	if !ok {
		panic("cant cast qlearning.State to game!")
	}
	gm.Move(s.Dir)
	return gm
}
