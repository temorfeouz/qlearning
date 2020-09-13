package main

var (
	wallTop = gameBlock{
		symbol:       "¯¯¯¯¯",
		isWin:        false,
		canGoThought: false,
	}
	wallBot = gameBlock{
		symbol:       "_____",
		isWin:        false,
		canGoThought: false,
	}
	wallLft = gameBlock{
		symbol:       "⎸",
		isWin:        false,
		canGoThought: false,
	}
	wallRgt = gameBlock{
		symbol:       "⎹",
		isWin:        false,
		canGoThought: false,
	}
	wallInr = gameBlock{
		symbol:       "|",
		isWin:        false,
		canGoThought: false,
	}
	blokEpt = gameBlock{
		symbol:       "\t",
		isWin:        false,
		canGoThought: true,
	}
	blokPrs = gameBlock{
		symbol:       "웃\t",
		isWin:        false,
		canGoThought: false,
	}
	blokWIN = gameBlock{
		symbol:       "♔",
		isWin:        true,
		canGoThought: true,
	}
)

type (
	gameBlock struct {
		symbol       string
		isWin        bool
		canGoThought bool
	}
	control struct {
		x int
		y int
	}
)

var (
	left = control{
		x: 0,
		y: -1,
	}
	right = control{
		x: 0,
		y: 1,
	}
	top = control{
		x: -1,
		y: 0,
	}
	bottom = control{
		x: 1,
		y: 0,
	}
)
