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
	files []string
	side  int // active list -- 0=dirs, 1=files

	visibleIdx  int // last visible file index on screen
	fileCount   int // number of files
	currentFile int // index of highllighted file

	visibleDirIdx int // last visible file index on screen
	dirCount      int // number of files
	currentDir    int // index of highllighted file

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

func showDirs(dirList []string) {

	for i, v := range dirList {

		for i >= visibleDirIdx && i < visibleDirIdx+(h-6) {
			if i == currentDir {
				fmt.Fprintf(os.Stdout, escapes.EraseLine)
				fmt.Println(bgCyan + brightWhite + v + reset)

				break
			} else {
				fmt.Fprintf(os.Stdout, escapes.EraseLine)
				fmt.Println(v)
				break
			}
		}
	}
}

func showFiles(dirList []string) {

	for i, v := range dirList {

		for i >= visibleDirIdx && i < visibleDirIdx+(h-6) {
			if i == currentDir {
				fmt.Fprintf(os.Stdout, escapes.EraseLine)
				fmt.Println(bgCyan + brightWhite + v + reset)
				break
			} else {
				fmt.Fprintf(os.Stdout, escapes.EraseLine)
				fmt.Println(v)
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
	visibleIdx = 0
	currentFile = 0
	visibleDirIdx = 0
	currentDir = 0
	dirs = nil
	files = nil
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
		}
		if ch == keyboard.KEY_RT {
			side = 1
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
			}
		}
	}

	// handle single key press
	// var ch int

	// for ch == 113 {
	// 	fmt.Fprintf(os.Stdout, "\033[?25h") // re-enable the cursor
	// 	os.Exit(0)
	// }
	// for ch != 113 {
	// 	ch = keyboard.ReadKey()

	// 	if visibleIdx >= 0 && currentFile > 0 && currentFile <= fileCount {
	// 		if ch == keyboard.KEY_UP {
	// 			currentFile--
	// 			if currentFile < visibleIdx {
	// 				visibleIdx--
	// 			}
	// 			listDirs()
	// 		}
	// 	}
	// 	if visibleIdx <= fileCount-1 && currentFile <= fileCount-2 {
	// 		if ch == keyboard.KEY_DN {
	// 			currentFile++
	// 			if currentFile > visibleIdx+(h-7) {
	// 				visibleIdx++
	// 			}
	// 			listDirs()
	// 		}
	// 	}
	// }
}

// func listFiles() {
// 	ansi.Theme("header", theme)

// 	// create a Slice containing paths from supplied Root
// 	getFiles, err := getAllFiles(root, "*.ans")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// iterate over each to split into 3 parts
// 	for i, v := range getFiles {

// 		dir, name := path.Split(v)
// 		base := filepath.Base(dir)

// 		f := FileData{
// 			Path: v,
// 			Base: base,
// 			Dir:  dir,
// 			Name: name,
// 		}

// 		for i >= visibleIdx && i < visibleIdx+(h-6) {
// 			if i == currentFile {
// 				fmt.Fprintf(os.Stdout, escapes.EraseLine)
// 				fmt.Println(bgCyan + brightWhite + f.FileInfo() + reset)
// 				break
// 			} else {
// 				fmt.Fprintf(os.Stdout, escapes.EraseLine)
// 				fmt.Println(f.FileInfo())
// 				break
// 			}
// 		}
// 	}
// }
// type FileData struct {
// 	Path string
// 	Dir  string
// 	Base string
// 	Name string
// }

// func (f FileData) FileInfo() string {
// 	// fmt.Fprintf(os.Stdout, f.Path)
// 	// fmt.Fprintf(os.Stdout, f.Dir)

// 	base := truncateText(f.Base, 35)
// 	name := truncateText(f.Name, 35)

// 	return escapes.CursorPosX(1) + base + escapes.CursorPosX(40) + name

// }

// func (f FileData) DirInfo() string {
// 	// fmt.Fprintf(os.Stdout, f.Path)
// 	// fmt.Fprintf(os.Stdout, f.Dir)

// 	base := truncateText(f.Base, 35)

// 	return escapes.CursorPosX(1) + base

// }

// func getAllFiles(root, pattern string) ([]string, error) {

// 	var matches []string
// 	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		if info.IsDir() {
// 			return nil
// 		}
// 		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
// 			return err
// 		} else if matched {
// 			matches = append(matches, path)
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	fileCount = len(matches)
// 	return matches, nil
// }
