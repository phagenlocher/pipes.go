package main

import (
	"flag"
	"fmt"
	"github.com/rthornton128/goncurses"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const UP = 0
const DOWN = 1
const RIGHT = 2
const LEFT = 3

var changeProb float64
var rand_start bool
var newColor bool
var waitTime time.Duration

func pipe(scr_lock chan bool) {
	// Generate color
	color := int16(rand.Intn(7) + 1)

	// Variables for direction
	dir := rand.Intn(3)
	var new_dir, old_dir int

	// Window and coordinates
	win := goncurses.StdScr()
	max_y, max_x := win.MaxYX()
	var x, y int
	if rand_start {
		x = rand.Intn(max_x)
		y = rand.Intn(max_y)
	} else {
		x = int(max_x / 2)
		y = int(max_y / 2)
	}

	for {
		// Store old directiion
		old_dir = dir
		if rand.Float64() > changeProb {
			// Get new direction
			new_dir = rand.Intn(4)
			// Check if the direction isn't the reversed
			// old direction.
			if ((new_dir + dir) % 4) != 1 {
				dir = new_dir
			}

		}

		// Get lock
		<-scr_lock
		// Set color
		win.ColorOn(color)
		// Print ACS char and change coordinates
		if dir == UP {
			if old_dir == LEFT {
				win.MoveAddChar(y, x, goncurses.ACS_LLCORNER)
			} else if old_dir == RIGHT {
				win.MoveAddChar(y, x, goncurses.ACS_LRCORNER)
			} else {
				win.MoveAddChar(y, x, goncurses.ACS_VLINE)
			}
			y--
		} else if dir == DOWN {
			if old_dir == LEFT {
				win.MoveAddChar(y, x, goncurses.ACS_ULCORNER)
			} else if old_dir == RIGHT {
				win.MoveAddChar(y, x, goncurses.ACS_URCORNER)
			} else {
				win.MoveAddChar(y, x, goncurses.ACS_VLINE)
			}
			y++
		} else if dir == RIGHT {
			if old_dir == UP {
				win.MoveAddChar(y, x, goncurses.ACS_ULCORNER)
			} else if old_dir == DOWN {
				win.MoveAddChar(y, x, goncurses.ACS_LLCORNER)
			} else {
				win.MoveAddChar(y, x, goncurses.ACS_HLINE)
			}
			x++
		} else if dir == LEFT {
			if old_dir == UP {
				win.MoveAddChar(y, x, goncurses.ACS_URCORNER)
			} else if old_dir == DOWN {
				win.MoveAddChar(y, x, goncurses.ACS_LRCORNER)
			} else {
				win.MoveAddChar(y, x, goncurses.ACS_HLINE)
			}
			x--
		}
		// Give back lock
		scr_lock <- true

		// Changing coordinates if leaving screen
		oob := true // Out of bounds
		if x > max_x {
			x = 0
		} else if y > max_y {
			y = 0
		} else if x < 0 {
			x = max_x
		} else if y < 0 {
			y = max_y
		} else {
			oob = false
		}
		// If the color needs to be changed and we went out of bounds
		// change the color
		if newColor && oob {
			color = int16(rand.Intn(7) + 1)
		}

		// Wait
		time.Sleep(waitTime)

	}

}

func main() {
	// Parse flags
	num_pipes := flag.Int("p", 1, "The `amount of pipes` to display")
	color := flag.Bool("C", false, "Disables color")
	NFlag := flag.Bool("N", false, "Changes the color of a pipe if it exits the screen")
	reset_lim := flag.Int("r", 2000, "Resets after the speciefied `amount of updates` (0 means no reset)")
	fps := flag.Int("f", 75, "Sets targeted `frames per second` (max 500)")
	sVal := flag.Float64("s", 0.8, "`Probability` of NOT changing the direction (0.0 - 1.0)")
	RFlag := flag.Bool("R", false, "Start at random coordinates")
	flag.Parse()

	// Set variables
	changeProb = *sVal
	rand_start = *RFlag
	newColor = *NFlag
	if *fps > 500 {
		waitTime = time.Duration(1000 / 500) * time.Millisecond
	} else if *fps != 0 {
		waitTime = time.Duration(1000 / *fps) * time.Millisecond
	} else {
		return
	}

	// Seeding RNG with current time
	rand.Seed(time.Now().Unix())

	// Disable SIGINT
	signal.Ignore(syscall.SIGINT)

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
	goncurses.CBreak(true)

	// Init colors
	goncurses.InitPair(1, goncurses.C_WHITE, goncurses.C_BLACK)
	goncurses.InitPair(2, goncurses.C_GREEN, goncurses.C_BLACK)
	goncurses.InitPair(3, goncurses.C_RED, goncurses.C_BLACK)
	goncurses.InitPair(4, goncurses.C_YELLOW, goncurses.C_BLACK)
	goncurses.InitPair(5, goncurses.C_BLUE, goncurses.C_BLACK)
	goncurses.InitPair(6, goncurses.C_MAGENTA, goncurses.C_BLACK)
	goncurses.InitPair(7, goncurses.C_CYAN, goncurses.C_BLACK)

	// Set timeout and clear
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
