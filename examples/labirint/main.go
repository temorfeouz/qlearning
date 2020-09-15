package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"go.uber.org/atomic"

	"github.com/temorfeouz/qlearning"
)

var lv = [][]gameBlock{
	{wallTop, wallTop, wallTop, wallTop, wallTop, wallTop, wallTop, wallTop, wallTop, wallTop},
	{wallLft, blokPrs, blokEpt, blokEpt, wallInr, blokEpt, blokEpt, blokEpt, blokEpt, wallRgt},
	{wallLft, blokEpt, blokEpt, blokEpt, wallInr, blokEpt, blokEpt, blokEpt, blokEpt, wallRgt},
	{wallLft, blokEpt, blokEpt, blokEpt, wallInr, blokEpt, blokEpt, blokEpt, blokEpt, wallRgt},
	{wallLft, blokEpt, blokEpt, blokEpt, blokEpt, blokEpt, blokEpt, blokEpt, blokEpt, wallRgt},
	{wallLft, blokEpt, blokEpt, blokEpt, blokEpt, blokEpt, blokEpt, blokEpt, blokEpt, wallRgt},
	{wallLft, blokEpt, blokEpt, blokEpt, wallInr, blokEpt, blokEpt, blokEpt, blokEpt, wallRgt},
	{wallLft, blokEpt, blokEpt, blokEpt, wallInr, blokEpt, blokEpt, blokEpt, blokEpt, wallRgt},
	{wallLft, blokEpt, blokEpt, blokEpt, wallInr, blokEpt, blokEpt, blokEpt, blokWIN, wallRgt},
	{wallBot, wallBot, wallBot, wallBot, wallBot, wallBot, wallBot, wallBot, wallBot, wallBot},
}

const qtableFile = "qtable.json"

var (
	autoplay bool = true
	steps    int  = 100
	rounds   int  = 10000
	threads       = runtime.NumCPU() * 2
)

func init() {
	flag.BoolVar(&autoplay, "autoplay", autoplay, "train")
	flag.IntVar(&steps, "steps", steps, "steps per round")
	flag.IntVar(&rounds, "rounds", rounds, "count of rounds")
	flag.IntVar(&threads, "threads", threads, "count of parallel games")
	flag.Parse()
}

func main() {
	if !autoplay {
		playByHangs()
	} else {
		train()
	}
}
func train() {
	agent := qlearning.NewSimpleAgent(0.7, 0.9)

	f, err := os.OpenFile(qtableFile, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	agent.Import(f)
	f.Close()

	var (
		round                  int
		refere                 = NewRefere(steps)
		wins, looses, minSteps int
		runedThreads           atomic.Int32
	)
	minSteps = steps // initial

	for round = 1; round <= rounds; round++ {
		runedThreads.Inc()
		go func() {
			gm := newGame(lv, true, false, round)

			done := false
			for {
				action := qlearning.Next(agent, gm, 0.7)
				agent.Learn(action, refere)

				win, st := gm.Stat()
				if win {
					wins++
					if st < minSteps {
						minSteps = st
					}
					done = true
				}
				if st >= steps {
					looses++
					done = true
				}

				gm.l("WINS:%d,LOOSES:%d,minSteps:%d, REW:%v", wins, looses, minSteps, refere.Reward(action))
				gm.Draw()
				//time.Sleep(50 * time.Millisecond)
				if done {
					runedThreads.Dec()
					return
				}
			}
		}()
		for int(runedThreads.Load()) >= threads {
			time.Sleep(time.Microsecond)
		}
	}

	f, err = os.OpenFile("qtable.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	agent.Export(f)
}

func playByHangs() {
	gm := newGame(lv, true, true, 1)

	gm.Draw()

	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	var b = make([]byte, 1)
	for {
		os.Stdin.Read(b)
		if b[0] == 97 {
			gm.l("move left")
			gm.Move(left)
		}
		if b[0] == 100 {
			gm.l("move right")
			gm.Move(right)
		}
		if b[0] == 119 {
			gm.l("move top")
			gm.Move(top)
		}
		if b[0] == 115 {
			gm.l("move bottom")
			gm.Move(bottom)
		}
		gm.LookAround()
		if isWin, steps := gm.Stat(); isWin {
			gm.l("WIN IN %d steps", steps)
		}

		gm.Draw()
	}
}
