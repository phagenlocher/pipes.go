package main

import (
	"flag"
	"fmt"
	"github.com/rthornton128/goncurses"
	"math/rand"
	"os"
	"time"
)

const UP = 0
const DOWN = 1
const RIGHT = 2
const LEFT = 3

var changeProb float64
var randStart bool
var newColor bool
var dimmedColors bool
var waitTime time.Duration

func pipe(screenLock chan bool) {
	// Generate color
	color := int16(rand.Intn(14) + 1)

	// Variables for curDirection
	curDir := rand.Intn(3)
	var newDir, oldDir int

	// Window and coordinates
	win := goncurses.StdScr()
	maxY, maxX := win.MaxYX()
	var x, y int
	if randStart {
		x = rand.Intn(maxX)
		y = rand.Intn(maxY)
	} else {
		x = int(maxX / 2)
		y = int(maxY / 2)
	}

	for {
		// Store old curDirectiion
		oldDir = curDir
		if rand.Float64() > changeProb {
			// Get new curDirection
			newDir = rand.Intn(4)
			// Check if the curDirection isn't the reversed
			// old curDirection.
			if ((newDir + curDir) % 4) != 1 {
				curDir = newDir
			}
		}

		// Generate color and dimming attribute
		dimmed := false
		nColor := color
		if color > 7 {
			dimmed = dimmedColors
			nColor -= 7
		}

		// Get lock
		<-screenLock
		// Set color and attribute
		if dimmed {
			win.AttrOn(goncurses.A_DIM)
		} else {
			win.AttrOff(goncurses.A_DIM)
		}
		win.ColorOn(nColor)
		// Print ACS char and change coordinates
		if curDir == UP {
			if oldDir == LEFT {
				win.MoveAddChar(y, x, goncurses.ACS_LLCORNER)
			} else if oldDir == RIGHT {
				win.MoveAddChar(y, x, goncurses.ACS_LRCORNER)
			} else {
				win.MoveAddChar(y, x, goncurses.ACS_VLINE)
			}
			y--
		} else if curDir == DOWN {
			if oldDir == LEFT {
				win.MoveAddChar(y, x, goncurses.ACS_ULCORNER)
			} else if oldDir == RIGHT {
				win.MoveAddChar(y, x, goncurses.ACS_URCORNER)
			} else {
				win.MoveAddChar(y, x, goncurses.ACS_VLINE)
			}
			y++
		} else if curDir == RIGHT {
			if oldDir == UP {
				win.MoveAddChar(y, x, goncurses.ACS_ULCORNER)
			} else if oldDir == DOWN {
				win.MoveAddChar(y, x, goncurses.ACS_LLCORNER)
			} else {
				win.MoveAddChar(y, x, goncurses.ACS_HLINE)
			}
			x++
		} else if curDir == LEFT {
			if oldDir == UP {
				win.MoveAddChar(y, x, goncurses.ACS_URCORNER)
			} else if oldDir == DOWN {
				win.MoveAddChar(y, x, goncurses.ACS_LRCORNER)
			} else {
				win.MoveAddChar(y, x, goncurses.ACS_HLINE)
			}
			x--
		}
		// Give back lock
		screenLock <- true

		// Changing coordinates if leaving screen
		oob := true // Out of bounds
		if x > maxX {
			x = 0
		} else if y > maxY {
			y = 0
		} else if x < 0 {
			x = maxX
		} else if y < 0 {
			y = maxY
		} else {
			oob = false
		}
		// If the color needs to be changed and we went out of bounds
		// change the color
		if newColor && oob {
			color = int16(rand.Intn(14) + 1)
		}

		// Wait
		time.Sleep(waitTime)

	}

}

func main() {
	// Parse flags
	num_pipes := flag.Int("p", 1, "The `amount of pipes` to display")
	color := flag.Bool("C", false, "Disables color")
	DFlag := flag.Bool("D", false, "Use dimmed colors in addition to normal colors")
	NFlag := flag.Bool("N", false, "Changes the color of a pipe if it exits the screen")
	reset_lim := flag.Int("r", 2000, "Resets after the speciefied `amount of updates` (0 means no reset)")
	fps := flag.Int("f", 75, "Sets targeted `frames per second` that also dictate the moving speed")
	sVal := flag.Float64("s", 0.8, "`Probability` of NOT changing the curDirection (0.0 - 1.0)")
	RFlag := flag.Bool("R", false, "Start at random coordinates")
	flag.Parse()

	// Set variables
	changeProb = *sVal
	randStart = *RFlag
	newColor = *NFlag
	dimmedColors = *DFlag
	// Set FPS
	if *fps > 1000000 {
		waitTime = time.Duration(1) * time.Microsecond
	} else if *fps > 0 {
		waitTime = time.Duration(1000000 / *fps) * time.Microsecond
	} else {
		// 0 or negative FPS are impossible
		return
	}

	// Seeding RNG with current time
	rand.Seed(time.Now().Unix())

	// Init ncurses
	stdscr, err := goncurses.Init()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer goncurses.End()

	// More init
	if !*color {
		goncurses.StartColor()
	}
	goncurses.FlushInput()
	goncurses.Cursor(0)
	goncurses.Echo(false)
	goncurses.Raw(true)

	// Init colors
	goncurses.UseDefaultColors()
	goncurses.InitPair(1, goncurses.C_WHITE, -1)
	goncurses.InitPair(2, goncurses.C_GREEN, -1)
	goncurses.InitPair(3, goncurses.C_RED, -1)
	goncurses.InitPair(4, goncurses.C_YELLOW, -1)
	goncurses.InitPair(5, goncurses.C_BLUE, -1)
	goncurses.InitPair(6, goncurses.C_MAGENTA, -1)
	goncurses.InitPair(7, goncurses.C_CYAN, -1)

	// Set timeout and clear
	stdscr.AttrSet(goncurses.A_NORMAL)
	stdscr.Timeout(0)
	stdscr.Clear()
	stdscr.Refresh()

	// Creat channel for lock
	lock := make(chan bool, 1)
	lock <- true

	// Generate goroutines
	for i := 0; i < *num_pipes; i++ {
		go pipe(lock)
	}

	// Refresh loop (runs until a key was pressed)
	for i := 0; stdscr.GetChar() == 0; {
		// Wait
		time.Sleep(waitTime)

		// Only increment if reset limited is not 0
		if *reset_lim != 0 {
			i++
		}

		// Reset limit has been reached
		if i > *reset_lim {
			stdscr.Clear()
			i = 0
		}
	}

}
