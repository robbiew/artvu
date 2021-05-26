package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/robbiew/artvu/ansi"
	term "github.com/robbiew/artvu/term"
	escapes "github.com/snugfox/ansi-escapes"
	keyboard "github.com/tlorens/go-ibgetkey"
)

var (
	h     int    // term height (rows)
	w     int    // term width (cols)
	root  string // art files dir
	theme int    // 80 or 132
	dirs  []string
	// files []string
	side int // active list -- 0=dirs, 1=files

	visibleFileIdx int // last visible file index on screen
	fileCount      int // number of files
	currentFile    int // index of highllighted file

	visibleDirIdx int // last visible file index on screen
	dirCount      int // number of files
	currentDir    int // index of highllighted file

	dirName string // name of current dir

	reset = "\u001b[0m"

	// colors
	black   = "\u001b[30m"
	red     = "\u001b[31m"
	green   = "\u001b[32m"
	yellow  = "\u001b[33m"
	blue    = "\u001b[34m"
	magenta = "\u001b[35m"
	cyan    = "\u001b[36m"
	white   = "\u001b[37m"

	brightBlack   = "\u001b[30;1m"
	brightRed     = "\u001b[31;1m"
	brightGreen   = "\u001b[32;1m"
	brightYellow  = "\u001b[33;1m"
	brightBlue    = "\u001b[34;1m"
	brightMagenta = "\u001b[35;1m"
	brightCyan    = "\u001b[36;1m"
	brightWhite   = "\u001b[37;1m"

	bgBlack   = "\u001b[40m"
	bgRed     = "\u001b[41m"
	bgGreen   = "\u001b[42m"
	bgYellow  = "\u001b[43m"
	bgBlue    = "\u001b[44m"
	bgMagenta = "\u001b[45m"
	bgCyan    = "\u001b[46m"
	bgWhite   = "\u001b[47m"
)

func createFilesSlice(root string, dir string) ([]string, error) {

	var files []string

	fileInfo, err := ioutil.ReadDir(root + "/" + dir)
	if err != nil {
		return files, err
	}
	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}

func createDirSlice(root string) ([]string, error) {

	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		dirs = append(dirs, file.Name())
	}
	dirCount = len(dirs)
	return dirs, nil
}

// Contains tells whether a contains x.
func contains(a []string, x string) int {
	var n int
	for _, n := range a {
		if x == n {
			return len(n)
		}
	}
	return n
}

func indexOf(data []string, element int) string {
	for k, v := range data {
		if element == k {
			return v
		}
	}
	return "not found" //not found.
}

func showFiles(files []string) {

	fmt.Fprintf(os.Stdout, "\033[6;32H")

	for i, v := range files {
		for i >= visibleFileIdx && i < visibleFileIdx+(h-6) {
			if i == currentFile {
				// active dir
				loc := "\033[32G"
				fmt.Println(loc + bgCyan + brightWhite + truncateText(v, 45) + reset)
				break
			} else {
				loc := "\033[32G"
				fmt.Println(loc + fmt.Sprint(i) + " " + truncateText(v, 45) + reset)

				break
			}

		}
	}
}

func showDirs(dirList []string) {

	for i, v := range dirList {
		for i >= visibleDirIdx && i < visibleDirIdx+(h-6) {
			if i == currentDir {
				// active dir

				fmt.Println(bgCyan + brightWhite + truncateText(v, 30) + reset)
				break
			} else {
				fmt.Println(truncateText(v, 30) + reset)
				break
			}
		}

	}
}

func truncateText(s string, max int) string {
	if len(s) > max {
		r := 0
		for i := range s {
			r++
			if r > max {
				return s[:i]
			}
		}
	}
	return s
}

func init() {
	visibleFileIdx = 0
	currentFile = 0
	visibleDirIdx = 0
	currentDir = 0
	dirs = nil
	side = 0
}

func main() {

	if len(os.Args) == 1 {
		log.Fatal("No art path given, Please specify path.")
		return
	}
	if root = os.Args[1]; root == "" {
		log.Fatal("No art path given, Please specify path.")
		return
	}

	fmt.Fprintf(os.Stdout, "\033[?25l") // hide the cursor
	fmt.Fprintf(os.Stdout, "\033[2J")   // clear the screen

	// Try and detect the user's term size
	h, w = term.GetTermSize()
	if w < 132 {
		theme = 80
	}
	if w >= 132 {
		theme = 132
	} else {
		theme = 80
	}

	fmt.Fprintf(os.Stdout, "\033[2J")   // clear the screen again
	fmt.Fprintf(os.Stdout, "\033[0;0f") // set cursor to 0,0 position

	// get and print list of top-level dirs under root
	list, err := createDirSlice(root)
	if err != nil {
		log.Fatal(err)
	}

	dirs = list
	ansi.Theme("header", theme)
	showDirs(dirs)

	// handle single key press
	var ch int
	for ch == 113 {
		fmt.Fprintf(os.Stdout, "\033[?25h") // re-enable the cursor
		os.Exit(0)
	}
	for ch != 113 {
		ch = keyboard.ReadKey()

		if ch == keyboard.KEY_LF {
			side = 0
			fmt.Fprintf(os.Stdout, "\033[2J")
			ansi.Theme("header", theme)

			fmt.Fprintf(os.Stdout, escapes.CursorPos(0, 5))
			showDirs(dirs)

		}
		if ch == keyboard.KEY_RT {
			side = 1

			currentDirName := indexOf(dirs, currentDir)
			filesSlice, err := createFilesSlice(root, currentDirName)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(os.Stdout, escapes.CursorPos(0, 5))
			showDirs(dirs)
			fmt.Fprintf(os.Stdout, escapes.CursorPos(0, 5))
			showFiles(filesSlice)

		}

		if ch == keyboard.KEY_UP {
			if side == 0 {
				if visibleDirIdx >= 0 && currentDir > 0 && currentDir <= dirCount {
					currentDir--
					if currentDir < visibleDirIdx {
						visibleDirIdx--
					}

					fmt.Fprintf(os.Stdout, escapes.CursorPos(0, 5))
					showDirs(dirs)
				}
			} else {
				if visibleFileIdx >= 0 && currentFile > 0 && currentFile <= fileCount {
					currentFile--
					if currentFile < visibleFileIdx {
						visibleFileIdx--
					}
					fmt.Fprintf(os.Stdout, escapes.CursorPos(0, 5))

				}
			}
		}
		if ch == keyboard.KEY_DN {
			if side == 0 {
				if visibleDirIdx <= dirCount-1 && currentDir <= dirCount-2 {
					currentDir++
					if currentDir > visibleDirIdx+(h-7) {
						visibleDirIdx++
					}
					fmt.Fprintf(os.Stdout, escapes.CursorPos(0, 5))
					showDirs(dirs)
				}
			} else {
				if visibleFileIdx <= fileCount-1 && currentFile <= fileCount-2 {
					currentFile++
					if currentFile > visibleFileIdx+(h-7) {
						visibleFileIdx++
					}
					fmt.Fprintf(os.Stdout, escapes.CursorPos(0, 5))

				}
			}
		}
	}
}
