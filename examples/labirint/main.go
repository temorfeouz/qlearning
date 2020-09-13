package main

import (
	"os"
	"os/exec"
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

func main() {
	gm := newGame(lv, true)

	gm.Draw()

	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	var b = make([]byte, 1)
	for {
		os.Stdin.Read(b)
		//fmt.Println("I got the byte", b, "("+string(b)+")")
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

		if isWin, steps := gm.Stat(); isWin {
			gm.l("WIN IN %d steps", steps)
		}

		gm.Draw()

	}
	// before entering the loop
	//fmt.Print("\033[s") // save the cursor position
	//fmt.Println("..")
	//i := 0

	//sigs := make(chan os.Signal, 1)
	//done := make(chan bool, 1)
	//signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	//for true {
	//	i++
	//	fmt.Print("\033[u\033[K") // restore the cursor position and clear the line
	//
	//	fmt.Print("\033[A") // move the cursor up
	//	fmt.Printf("Retrieved %d\n", i)
	//	time.Sleep(time.Second)
	//
	//}
	//<-done
	//fmt.Println("EXIT")
}

type gamerable interface {
	Draw()
}

func loop(gm gamerable, control chan control) {
	// process control
	for {
		select {
		//case :

		}
		// draw level
	}

}
