# pipes.go

![Screenshot](doc/screen.png "Screenshot")

## Description
A [pipes.sh](https://github.com/pipeseroni/pipes.sh) clone written in Golang. While the original cannot be beaten in terms of compatibility and file size, it has a bit of a performance issue.
*pipes.go* tries to fix this problem with high concurrency. This results in lower CPU usage and smoother animation.

## Requirements
* Golang
* ncurses

## Usage
* **-C** disables color
* **-B** disables bold output
* **-D** additionally uses dimmed colors
* **-N** lets the pipes change color when exiting the screen (just like in pipes.sh)
* **-R** lets the pipes start from random coordinates
* **-p** specifies the amount of pipes
* **-c** sets a predefined colorscheme
* **-t** sets the character set
* **-r** specifies after how many updates to clear the screen
* **-f** sets the targeted frames per second
* **-s** sets the probability of not changing the direction for a pipe

## Building & Installation
```
make
```
A simple make call will install the needed library [goncurses](https://github.com/rthornton128/goncurses) and compile the program.
```
make install
```
This will build the executable and move it to */usr/local/bin*.
```
make uninstall
```
This will delete the executable from */usr/local/bin*.
