package main

var (
	wallTop = gameBlock{
		symbol:       "▁▁",
		isWin:        false,
		canGoThought: false,
	}
	wallBot = gameBlock{
		symbol:       "▔▔",
		isWin:        false,
		canGoThought: false,
	}
	wallLft = gameBlock{
		symbol:       "▊",
		isWin:        false,
		canGoThought: false,
	}
	wallRgt = gameBlock{
		symbol:       "▊",
		isWin:        false,
		canGoThought: false,
	}
	wallInr = gameBlock{
		symbol:       "▊░",
		isWin:        false,
		canGoThought: false,
	}
	blokEpt = gameBlock{
		symbol:       "░░",
		isWin:        false,
		canGoThought: true,
	}
	blokPrs = gameBlock{
		symbol:       "*░",
		isWin:        false,
		canGoThought: false,
	}
	blokWIN = gameBlock{
		symbol:       "♔ ",
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
		x   int
		y   int
		dir rune
	}
)

func (c *control) String() string {
	return string(c.dir)
}

var (
	left = control{
		dir: '←',
		x:   0,
		y:   -1,
	}
	right = control{
		dir: '→',
		x:   0,
		y:   1,
	}
	top = control{
		dir: '↑',
		x:   -1,
		y:   0,
	}
	bottom = control{
		dir: '↓',
		x:   1,
		y:   0,
	}
)

func (b *gameBlock) String() string {
	return b.symbol
}
