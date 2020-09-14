package main

import (
	"fmt"

	"github.com/temorfeouz/qlearning"
)

type Step struct {
	Dir         control
	BlockAtStep gameBlock
	DistToWin   float64
}

func (s *Step) String() string {
	return fmt.Sprintf("cs:%s:d:%v", s.BlockAtStep.String(), s.DistToWin)
}
func (s *Step) Apply(st qlearning.State) qlearning.State {
	gm, ok := st.(*game)
	if !ok {
		panic("cant cast qlearning.State to game!")
	}
	gm.Move(s.Dir)
	return gm
}
