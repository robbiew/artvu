package main

import (
	"fmt"
	"os"

	term "github.com/robbiew/artvu-ansi-gallery/term"
	keyboard "github.com/tlorens/go-ibgetkey"
)

var (
	h int // term height (rows)
	w int // term width (cols)
)

func main() {

	fmt.Fprintf(os.Stdout, "\033[2J") // clear the screen
	h, w = term.GetTermSize()
	fmt.Fprintf(os.Stdout, "\033[2J")   // clear the screen again
	fmt.Fprintf(os.Stdout, "\033[0;0f") // set cursor to 0,0 position

	fmt.Println(h, w)

	// handle single key press
	var ch int

	for ch == 113 {
		os.Exit(0)
	}
	for ch != 113 {
		ch = keyboard.ReadKey()
		if ch == keyboard.KEY_RT {
		}
		if ch == keyboard.KEY_LF {
		}
	}
}
