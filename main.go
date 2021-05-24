package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/robbiew/artvu/ansi"
	term "github.com/robbiew/artvu/term"
	keyboard "github.com/tlorens/go-ibgetkey"
)

var (
	h     int    // term height (rows)
	w     int    // term width (cols)
	root  string // art files dir
	theme int    // 80 or 132

	visibleIdx  int // last visible file index on screen
	fileCount   int // number of files
	currentFile int // index of highllighted file
)

type FileData struct {
	Path string
	Dir  string
	Base string
	Name string
}

func (f FileData) FileInfo() string {
	// fmt.Fprintf(os.Stdout, f.Path)
	// fmt.Fprintf(os.Stdout, f.Dir)
	return truncateText(f.Base, 35) + truncateText(f.Name, 35)

}

func getAllFiles(root, pattern string) ([]string, error) {

	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	fileCount = len(matches)
	// fmt.Println(fileCount)

	return matches, nil
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

func listFiles() {

	// create a Slice containing paths from supplied Root
	getFiles, err := getAllFiles(root, "*.ans")
	if err != nil {
		panic(err)
	}

	// iterate over each to split into 3 parts

	for i, v := range getFiles {

		dir, name := path.Split(v)
		base := filepath.Base(dir)

		f := FileData{
			Path: v,
			Base: base,
			Dir:  dir,
			Name: name,
		}

		for i > visibleIdx && i < visibleIdx+(h-3) {
			fmt.Println(i, f.FileInfo())
			break
		}

	}

}

func init() {
	visibleIdx = 0
	currentFile = 0
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

	fmt.Fprintf(os.Stdout, "\033[2J") // clear the screen
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

	ansi.Theme("header", theme)
	fmt.Println("\r")

	listFiles()

	// handle single key press
	var ch int

	for ch == 113 {
		fmt.Fprintf(os.Stdout, "\033[?25h") // re-enable the cursor
		os.Exit(0)
	}
	for ch != 113 {
		ch = keyboard.ReadKey()

		if visibleIdx >= 0 && currentFile > 0 {
			if ch == keyboard.KEY_UP {
				currentFile--
				if currentFile <= visibleIdx {
					visibleIdx--
				}
				fmt.Fprintf(os.Stdout, "\033[4;0f")
				fmt.Fprintf(os.Stdout, "\033[2J")
				listFiles()
				fmt.Fprintf(os.Stdout, "\033[0;0f")
				fmt.Println(currentFile)
			}
		}
		if visibleIdx < fileCount-(h-3) && currentFile < fileCount {
			if ch == keyboard.KEY_DN {
				currentFile++
				if currentFile > visibleIdx+(h-4) {
					visibleIdx++
				}
				fmt.Fprintf(os.Stdout, "\033[4;0f")
				fmt.Fprintf(os.Stdout, "\033[2J")
				listFiles()
				fmt.Fprintf(os.Stdout, "\033[0;0f")
				fmt.Println(currentFile)
			}
		}
	}
}
