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

func pipe(x, y int, ch_prob float64, scr_lock chan bool) {
	// Generate color
	color := int16(rand.Intn(7) + 1)

	// Variables for direction
	dir := rand.Intn(3)
	var new_dir, old_dir int

	// Window and coordinates
	win := goncurses.StdScr()
	max_y, max_x := win.MaxYX()

	for {

		// Store old directiion
		old_dir = dir
		if rand.Float64() > ch_prob {
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
		// Give back log
		scr_lock <- true

		// Changing coordinates if leaving screen
		if x > max_x {
			x = 1
		}
		if y > max_y {
			y = 1
		}
		if x < 1 {
			x = max_x
		}
		if y < 1 {
			y = max_y
		}

		// Wait
		time.Sleep(20 * time.Millisecond)

	}

}

func main() {
	// Parse flags
	num_pipes := flag.Int("p", 1, "The `amount of pipes` to display")
	color := flag.Bool("C", false, "Disables color")
	reset_lim := flag.Int("r", 2000, "Resets after the speciefied `amount of updates`")
	ch_prob := flag.Float64("s", 0.8, "`Probability` of NOT changing the direction (0.0 - 1.0)")
	random := flag.Bool("R", false, "Start at random coordinates")
	flag.Parse()

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
	goncurses.Echo(true)
	goncurses.CBreak(false)

	// Init colors
	goncurses.InitPair(1, goncurses.C_WHITE, goncurses.C_BLACK)
	goncurses.InitPair(2, goncurses.C_GREEN, goncurses.C_BLACK)
	goncurses.InitPair(3, goncurses.C_RED, goncurses.C_BLACK)
	goncurses.InitPair(4, goncurses.C_YELLOW, goncurses.C_BLACK)
	goncurses.InitPair(5, goncurses.C_BLUE, goncurses.C_BLACK)
	goncurses.InitPair(6, goncurses.C_MAGENTA, goncurses.C_BLACK)
	goncurses.InitPair(7, goncurses.C_CYAN, goncurses.C_BLACK)

	// Set timeout and get max coordinates
	stdscr.Timeout(0)
	max_y, max_x := stdscr.MaxYX()

	// Creat channel for lock
	lock := make(chan bool, 1)
	lock <- true

	// Generate goroutines (the pipes :P)
	for i := 0; i < *num_pipes; i++ {
		if *random {
			go pipe(rand.Intn(max_x), rand.Intn(max_y), *ch_prob, lock)
		} else {
			go pipe(int(max_x/2), int(max_y/2), *ch_prob, lock)
		}
	}

	// Refresh loop
	for i := 0; stdscr.GetChar() == 0; i++ {
		time.Sleep(10 * time.Millisecond)

		if i > *reset_lim {
			stdscr.Clear()
			i = 0
		}
	}

}
