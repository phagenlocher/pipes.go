package main

import (
    "flag"
    "fmt"
    "github.com/rthornton128/goncurses"
    "math/rand"
    "os"
    "time"
)

const (
    UP    = 0
    DOWN  = 1
    RIGHT = 2
    LEFT  = 3
)

var printChars [6]goncurses.Char
var changeProb float64
var randStart bool
var newColor bool
var dimmedColors bool
var numColors int
var waitTime time.Duration

func pipe(screenLock chan bool) {
    // Generate color
    color := int16(rand.Intn(numColors * 2) + 1)

    // Character to be printed
    var printChar goncurses.Char

    // Variables for dircetion
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
        // Store old direction
        oldDir = curDir
        if rand.Float64() > changeProb {
            // Get new direction
            newDir = rand.Intn(4)
            // Check if the direction isn't the reversed
            // old direction.
            if ((newDir + curDir) % 4) != 1 {
                curDir = newDir
            }
        }

        // Generate color and dimming attribute
        dimmed := false
        nColor := color
        if color > int16(numColors) {
            // Only dim if the flag has been set
            dimmed = dimmedColors
            // Subtract num of colors to get actual color
            nColor -= int16(numColors)
        }

        // Get char corresponding to pipes direction
        if curDir == UP {
            if oldDir == LEFT {
                printChar = printChars[4]
            } else if oldDir == RIGHT {
                printChar = printChars[5]
            } else {
                printChar = printChars[1]
            }
        } else if curDir == DOWN {
            if oldDir == LEFT {
                printChar = printChars[2]
            } else if oldDir == RIGHT {
                printChar = printChars[3]
            } else {
                printChar = printChars[1]
            }
        } else if curDir == RIGHT {
            if oldDir == UP {
                printChar = printChars[2]
            } else if oldDir == DOWN {
                printChar = printChars[4]
            } else {
                printChar = printChars[0]
            }
        } else if curDir == LEFT {
            if oldDir == UP {
                printChar = printChars[3]
            } else if oldDir == DOWN {
                printChar = printChars[5]
            } else {
                printChar = printChars[0]
            }
        }

        // Get lock
        <-screenLock
        // Set attribute
        if dimmed {
            win.AttrOn(goncurses.A_DIM)
        } else {
            win.AttrOff(goncurses.A_DIM)
        }
        // Set color
        win.ColorOn(nColor)
        // Print char
        win.MoveAddChar(y, x, printChar)
        // Give back lock
        screenLock <- true

        // Update coordinates
        if curDir == UP {
            y--
        } else if curDir == DOWN {
            y++
        } else if curDir == RIGHT {
            x++
        } else if curDir == LEFT {
            x--
        }

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
            color = int16(rand.Intn(numColors * 2) + 1)
        }

        // Wait
        time.Sleep(waitTime)

    }

}

func setPrintChars(set int) {
    switch set {
    default:
        printChars[0] = goncurses.ACS_HLINE
        printChars[1] = goncurses.ACS_VLINE
        printChars[2] = goncurses.ACS_ULCORNER
        printChars[3] = goncurses.ACS_URCORNER
        printChars[4] = goncurses.ACS_LLCORNER
        printChars[5] = goncurses.ACS_LRCORNER
    case 1:
        printChars[0] = '.'
        printChars[1] = '.'
        printChars[2] = 'o'
        printChars[3] = 'o'
        printChars[4] = 'o'
        printChars[5] = 'o'
    case 2:
        printChars[0] = '.'
        printChars[1] = '.'
        printChars[2] = '.'
        printChars[3] = '.'
        printChars[4] = '.'
        printChars[5] = '.'
    case 3:
        printChars[0] = '-'
        printChars[1] = '|'
        printChars[2] = '+'
        printChars[3] = '+'
        printChars[4] = '+'
        printChars[5] = '+'
    }
}

func setColorScheme(scheme int) int {
    // Try to use the default background
    var background int16
    err := goncurses.UseDefaultColors()
    if err != nil {
        background = goncurses.C_BLACK
    } else {
        background = -1
    }

    // Init pairs according to scheme
    switch scheme {
    default:
        goncurses.InitPair(1, goncurses.C_WHITE,   background)
        goncurses.InitPair(2, goncurses.C_GREEN,   background)
        goncurses.InitPair(3, goncurses.C_RED,     background)
        goncurses.InitPair(4, goncurses.C_YELLOW,  background)
        goncurses.InitPair(5, goncurses.C_BLUE,    background)
        goncurses.InitPair(6, goncurses.C_MAGENTA, background)
        goncurses.InitPair(7, goncurses.C_CYAN,    background)
        return 7
    case 1:
        goncurses.InitPair(1, goncurses.C_WHITE, background)
        goncurses.InitPair(2, goncurses.C_BLUE,  background)
        goncurses.InitPair(3, goncurses.C_CYAN,  background)
        return 3
    case 2:
        goncurses.InitPair(1, goncurses.C_RED,    background)
        goncurses.InitPair(2, goncurses.C_YELLOW, background)
        goncurses.InitPair(3, goncurses.C_GREEN,  background)
        return 3
    case 3:
        goncurses.InitPair(1, goncurses.C_WHITE, background)
        goncurses.InitPair(2, goncurses.C_BLUE,  background)
        goncurses.InitPair(3, goncurses.C_RED,   background)
        return 3
    case 4:
        goncurses.InitPair(1, goncurses.C_RED,   background)
        goncurses.InitPair(2, goncurses.C_GREEN, background)
        goncurses.InitPair(3, goncurses.C_BLUE,  background)
        return 3
    case 5:
        goncurses.InitPair(1, goncurses.C_WHITE, background)
        goncurses.InitPair(2, goncurses.C_RED,   background)
        return 2
    case 6:
        goncurses.InitPair(1, goncurses.C_WHITE, background)
        goncurses.InitPair(2, goncurses.C_BLUE,  background)
        return 2
    case 7:
        goncurses.InitPair(1, goncurses.C_WHITE, background)
        goncurses.InitPair(2, goncurses.C_GREEN, background)
        return 2
    }
}

func main() {
    // Parse flags
    color       := flag.Bool("C", false, "Disables color")
    BFlag       := flag.Bool("B", false, "Disables bold output")
    RFlag       := flag.Bool("R", false, "Start at random coordinates")
    DFlag       := flag.Bool("D", false, "Use dimmed colors in addition to normal colors (impacts performance drastically)")
    NFlag       := flag.Bool("N", false, "Changes the color of a pipe if it exits the screen")
    numPipes    := flag.Int("p", 5, "The `amount of pipes` to display")
    resetLim    := flag.Int("r", 600, "Resets after the speciefied `amount of updates` (0 means no reset)")
    fps         := flag.Int("f", 60, "Sets targeted `frames per second` that also dictate the moving speed")
    colorScheme := flag.Int("c", 0, "Sets the `colorscheme` (0-7)")
    charSet     := flag.Int("t", 0, "Sets the `character set` (0-3)")
    sVal        := flag.Float64("s", 0.8, "`Probability` of NOT changing the direction (0.0-1.0)")
    flag.Parse()

    // Set variables
    changeProb   = *sVal
    randStart    = *RFlag
    newColor     = *NFlag
    dimmedColors = *DFlag
    // Set FPS
    if *fps > 1000000 {
        waitTime = time.Duration(1) * time.Microsecond
    } else if *fps > 0 {
        waitTime = time.Duration(1000000 / *fps) * time.Microsecond
    } else {
        // 0 or negative FPS are impossible
        fmt.Fprintln(os.Stderr, "FPS cannot be smaller than 1")
        os.Exit(1)
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

    // Init color pairs and number of colors
    numColors = setColorScheme(*colorScheme)

    // Init chars
    setPrintChars(*charSet)

    // Set attribute, timeout (non blocking) and clear screen
    if *BFlag {
        stdscr.AttrSet(goncurses.A_NORMAL)
    } else {
        stdscr.AttrSet(goncurses.A_BOLD)
    }
    stdscr.Timeout(0)
    stdscr.Clear()
    stdscr.Refresh()

    // Create channel for lock
    lock := make(chan bool, 1)
    lock <- true

    // Generate goroutines
    for i := 0; i < *numPipes; i++ {
        go pipe(lock)
    }

    // Refresh loop (runs until a key was pressed)
    for i := 0; stdscr.GetChar() == 0; {
        // Wait
        time.Sleep(waitTime)

        // Reset limit has been reached
        if i > *resetLim {
            stdscr.Clear()
            i = 0
        } else if *resetLim != 0 {
            // Only increment if reset limit is not 0
            i++
        }
    }

}
