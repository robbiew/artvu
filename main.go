package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/robbiew/artvu/sauce"
	term "github.com/robbiew/artvu/term"
	escapes "github.com/snugfox/ansi-escapes"
	keyboard "github.com/tlorens/go-ibgetkey"
)

var (
	h       int      // term height (rows)
	w       int      // term width (cols)
	headerH int      // height of header
	root    string   // art files dir
	theme   int      // 80 or 132
	dirs    []string // dir list
	side    int      // active list -- 0=dirs, 1=files
	canQuit bool

	splitCol int // where the second column should begin

	visibleFileIdx int // last visible file index on screen
	fileCount      int // number of files
	currentFile    int // index of highllighted file

	visibleDirIdx   int // last visible file index on screen
	dirCount        int // number of files
	currentDir      int // index of highllighted file
	currentFileName string

	fileHasSAUCE bool

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

type WalkFunc func(path string, info os.FileInfo, err error) error

func Theme(name string, size int) {

	s := strconv.Itoa(size)
	file := "theme/" + name + "." + s + ".ans"
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string and print to screen
	text := string(content)
	trimmed := TrimStringFromSauce(text)
	fmt.Fprintf(os.Stdout, trimmed)

	return
}

func TrimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}

func TrimStringFromSauce(s string) string {
	if idx := strings.Index(s, "SAUCE00"); idx != -1 {
		string := s
		delimiter := "SAUCE00"
		leftOfDelimiter := strings.Split(string, delimiter)[0]
		trim := TrimLastChar(leftOfDelimiter)
		return trim
		// rightOfDelimiter := strings.Join(strings.Split(string, delimiter)[1:], delimiter)
	}
	return s
}

// showArt draws the art
func ShowArt(path string, sizeX int, sizeY int, delay int) {

	rowCount := 0

	file, err := os.Open(path)
	if err != nil {
		//handle error
		return
	}

	defer file.Close()
	s := bufio.NewScanner(file)

	defer file.Close()
	fmt.Fprintf(os.Stdout, escapes.EraseScreen)
	// fmt.Println(escapes.CursorPos(0, 0))

	for s.Scan() {
		read_line := s.Text()
		// trim the text if it's after a SAUCE RECORD
		trimmed := TrimStringFromSauce(read_line)
		var b bytes.Buffer
		for {
			// add delay between each line to throttle speed
			fmt.Println(escapes.CursorPos(0, rowCount))
			time.Sleep(time.Duration(delay) * time.Millisecond)
			// fmt.Fprintf(os.Stdout, escapes.CursorNextLine)
			b.Write([]byte(trimmed + "\r"))
			b.WriteTo(os.Stdout)
			rowCount++
			break
		}
	}

}

func checkSauce(input string) bool {
	// let's check the file for a valid SAUCE record
	record := sauce.GetSauce(input)

	// if we find a SAUCE record, update bool flag
	if string(record.Sauceinf.ID[:]) == sauce.SauceID {
		fileHasSAUCE = true
		return true
	} else {
		return false
	}
}

func createFilesSlice(root string, dir string) ([]string, error) {

	rootDir := root + "/" + dir

	var files []string
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".ans" {
			files = append(files, info.Name())
		}
		return nil
	})
	return files, err
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

func indexOf(data []string, element int) string {
	for k, v := range data {
		if element == k {
			return v
		}
	}
	return "not found" //not found.
}

func scrollArrows(side int) {

	if side == 1 {
		// up arrow
		scrollUpLeft := "\033[" + fmt.Sprint(headerH+1) + ";" + fmt.Sprint(splitCol) + "H"
		fmt.Fprintf(os.Stdout, scrollUpLeft)
		fmt.Fprintf(os.Stdout, string([]rune{'\u0018'}))

		//down arrow
		scrollDownLeft := "\033[" + fmt.Sprint(headerH+2) + ";" + fmt.Sprint(splitCol) + "H"
		fmt.Fprintf(os.Stdout, scrollDownLeft)
		fmt.Fprintf(os.Stdout, string([]rune{'\u0019'}))

	}

	if side == 0 {
		// up arrow
		scrollUpRight := "\033[" + fmt.Sprint(headerH+1) + ";2H"
		fmt.Fprintf(os.Stdout, scrollUpRight)
		fmt.Fprintf(os.Stdout, string([]rune{'\u0018'}))

		//down arrow
		scrollDownRight := "\033[" + fmt.Sprint(headerH+2) + ";2H"
		fmt.Fprintf(os.Stdout, scrollDownRight)
		fmt.Fprintf(os.Stdout, string([]rune{'\u0019'}))
	}

}

func showWidth(files []string) {

	widthLoc := "\033[" + fmt.Sprint(headerH+1) + ";" + fmt.Sprint(splitCol-4) + "H"
	fmt.Fprintf(os.Stdout, widthLoc)

	for i, v := range files {
		for i >= visibleFileIdx && i < visibleFileIdx+(h-(headerH+2)) {

			// active dir
			loc := "\033[" + fmt.Sprint(splitCol-4) + "G"
			currentFileName = v
			currentDirName := indexOf(dirs, currentDir)
			file := root + "/" + currentDirName + "/" + currentFileName

			// Does the file have a SAUCE record?
			x := checkSauce(file)
			var color string
			if x == true {
				// Get width from Sauce
				allRecords := sauce.GetSauce(file)
				wi := allRecords.Sauceinf.Tinfo1
				if (int(wi)) <= w {
					color = green
				} else {
					color = red
				}
				if wi != 0 {
					fmt.Println(loc + color + strconv.Itoa(int(wi)) + reset)
				} else {
					fmt.Println(loc + red + "???" + reset)
				}
				break
			}
		}
	}
}

func showFiles(files []string) {

	// simulated scroll bar
	// arrowHeight := float64(1)
	// viewportHeight := float64(h - 6)                         // 19
	// contentHeight := float64(len(files))                     // 100
	// viewableRatio := float64(viewportHeight / contentHeight) // .19
	// scrollBarArea := viewportHeight - arrowHeight            // 18
	// thumbHeight := math.Round(scrollBarArea * viewableRatio) //3.42

	//

	filesLoc := "\033[" + fmt.Sprint(headerH+1) + ";" + fmt.Sprint(splitCol+1) + "H"
	fmt.Fprintf(os.Stdout, filesLoc)

	for i, v := range files {
		for i >= visibleFileIdx && i < visibleFileIdx+(h-(headerH+2)) {
			if i == currentFile {
				// active dir
				currentFileName = v
				loc := "\033[" + fmt.Sprint(splitCol+1) + "G"
				up := "\033[1A" // move cursor up
				fmt.Println(loc + " " + bgCyan + PadLeft(">", " ", splitCol-2) + reset)
				fmt.Println(up + loc + " " + reset + bgCyan + brightWhite + " " + truncateText(v, splitCol-4) + " " + reset)
				break
			} else {
				loc := "\033[" + fmt.Sprint(splitCol+1) + "G"
				fmt.Println(reset + loc + " " + brightBlack + " " + truncateText(v, splitCol-4) + reset)
				break
			}

		}
	}

	// simulated scroll bar
	// if isMultipleOf(thumbHeight, currentFile) == true {
	// 	scrollBar := "\033[" + fmt.Sprint(currentFile+7) + ";" + fmt.Sprint(80) + "H"
	// 	fmt.Fprintf(os.Stdout, scrollBar)
	// 	fmt.Fprintf(os.Stdout, bgCyan+" "+reset)
	// }
	scrollArrows(1)
}

func showDirs(dirList []string) {

	dirsLoc := "\033[" + fmt.Sprint(headerH+1) + ";" + fmt.Sprint(0) + "H"
	fmt.Fprintf(os.Stdout, dirsLoc)

	for i, v := range dirList {
		for i >= visibleDirIdx && i < visibleDirIdx+(h-(headerH+1)) {
			if i == currentDir {
				// active dir
				up := "\033[1A" // move cursor up
				fmt.Println(" " + bgCyan + "  " + PadLeft(">", " ", splitCol-10) + reset)
				fmt.Println(up + " " + reset + "  " + bgCyan + brightWhite + " " + truncateText(v, splitCol-11) + " " + reset)

				break
			} else {
				up := "\033[1A" // move cursor up
				fmt.Println(" " + cyan + "  " + PadLeft(">", " ", splitCol-8) + reset)
				fmt.Println(up + " " + reset + "  " + cyan + " " + truncateText(v, splitCol-11) + reset)
				break
			}
		}

	}
	scrollArrows(0)
}

func PadRight(str, pad string, lenght int) string {
	for {
		str += pad
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}

func PadLeft(str, pad string, lenght int) string {
	for {
		str = pad + str
		if len(str) > lenght {
			return str[0:lenght]
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
	headerH = 4
	canQuit = false
}

func showDiz(root string, dir string) {

	file := root + "/" + dir + "/file_id.diz"
	content, err := ioutil.ReadFile(file)

	if err != nil {
		log.Fatal(err)
	}

	rowCount := headerH + 1

	s := bufio.NewScanner(strings.NewReader(string(content)))

	// defer file.Close()
	for s.Scan() {
		read_line := s.Text()
		// trim the text if it's after a SAUCE RECORD
		trimmed := TrimStringFromSauce(read_line)
		loc := "\033[" + fmt.Sprint(rowCount) + ";" + fmt.Sprint(37) + "f"
		var b bytes.Buffer
		for rowCount < h-2 {
			b.Write([]byte(loc + trimmed + "\r"))
			b.WriteTo(os.Stdout)
			rowCount++
			break
		}

	}

}

func isMultipleOf(thumbHeight float64, rowNum int) bool {

	r := thumbHeight
	if math.Mod(float64(rowNum), r) == 0 {
		return true
	} else {
		return false
	}

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
		splitCol = 40
	}
	if w >= 132 {
		theme = 132
		splitCol = 40
	} else {
		theme = 80
		splitCol = 40
	}

	fmt.Fprintf(os.Stdout, "\033[2J")   // clear the screen again
	fmt.Fprintf(os.Stdout, "\033[0;0f") // set cursor to 0,0 position

	// get and print list of top-level dirs under root
	list, err := createDirSlice(root)
	if err != nil {
		log.Fatal(err)
	}

	dirs = list
	Theme("header", theme)
	showDirs(dirs)

	bottom := "\033[" + fmt.Sprint(h) + ";0H"
	fmt.Fprintf(os.Stdout, bottom)
	Theme("footer", theme)
	side = 0
	canQuit = false

	// handle single key press
	var ch int

	// for ch == 113 { // Quit

	// 	if side == 0 && canQuit == true {
	// 		fmt.Fprintf(os.Stdout, "\033[?25h") // re-enable the cursor
	// 		os.Exit(0)
	// 	}
	// 	if side == 1 && canQuit == false {
	// 		fmt.Println("q pressed")
	// 	} else {

	// 	}
	// }

	for ch != 113 {

		ch = keyboard.ReadKey()

		if ch == keyboard.KEY_LF { // LEFT ARROW
			if side == 1 {
				side = 0
				fmt.Fprintf(os.Stdout, "\033[2J")
				Theme("header", theme)
				showDirs(dirs)
				bottom := "\033[" + fmt.Sprint(h) + ";0H"
				fmt.Fprintf(os.Stdout, bottom)
				Theme("footer", theme)
			} else {

			}
		}

		if ch == keyboard.KEY_RT || ch == 10 { // RIGHT ARROW
			if side == 0 {
				side = 1
				currentFile = 0
				visibleFileIdx = 0
				currentDirName := indexOf(dirs, currentDir)
				filesSlice, err := createFilesSlice(root, currentDirName)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Fprintf(os.Stdout, "\033[2J") //clear screen
				Theme("header", theme)            // draw header

				showFiles(filesSlice)
				showDirs(dirs)
				showWidth(filesSlice)

				bottom := "\033[" + fmt.Sprint(h) + ";0H"
				fmt.Fprintf(os.Stdout, bottom)
				Theme("footer", theme)
			}
		}

		if ch == 10 { // Return
			if side == 1 {
				currentDirName := indexOf(dirs, currentDir)
				file := root + "/" + currentDirName + "/" + currentFileName
				ShowArt(file, w, h, 70)
			}
			if side == 0 {
				fmt.Println("enter pressed")

			}
		}

		if ch == keyboard.KEY_UP { // SCROLL UP
			if side == 0 {
				if visibleDirIdx >= 0 && currentDir > 0 && currentDir <= dirCount {
					currentDir--
					if currentDir < visibleDirIdx {
						visibleDirIdx--
					}

					fmt.Fprintf(os.Stdout, "\033[2J")
					Theme("header", theme)

					showDirs(dirs)

					bottom := "\033[" + fmt.Sprint(h) + ";0H"
					fmt.Fprintf(os.Stdout, bottom)
					Theme("footer", theme)
				}
			}

			if side == 1 {
				if visibleFileIdx >= 0 && currentFile > 0 && currentFile <= fileCount {
					currentFile--
					if currentFile < visibleFileIdx {
						visibleFileIdx--
					}

					currentDirName := indexOf(dirs, currentDir)

					filesSlice, err := createFilesSlice(root, currentDirName)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Fprintf(os.Stdout, "\033[2J")
					Theme("header", theme)

					showWidth(filesSlice)
					showFiles(filesSlice)
					showDirs(dirs)

					bottom := "\033[" + fmt.Sprint(h) + ";0H"
					fmt.Fprintf(os.Stdout, bottom)
					Theme("footer", theme)
				}
			}
		}

		if ch == keyboard.KEY_DN { //SCROLL DOWN
			if side == 0 {
				if visibleDirIdx <= dirCount-1 && currentDir <= dirCount-2 {
					currentDir++
					if currentDir > visibleDirIdx+(h-(headerH+3)) {
						visibleDirIdx++
					}
					showDirs(dirs)
					bottom := "\033[" + fmt.Sprint(h) + ";0H"
					fmt.Fprintf(os.Stdout, bottom)
					Theme("footer", theme)
				}
			}

			if side == 1 {
				currentDirName := indexOf(dirs, currentDir)
				filesSlice, err := createFilesSlice(root, currentDirName)
				if err != nil {
					log.Fatal(err)
				}
				fileCount = len(filesSlice)

				if visibleFileIdx <= fileCount-1 && currentFile <= fileCount-2 {
					currentFile++
					if currentFile > visibleFileIdx+(h-(headerH+3)) {
						visibleFileIdx++
					}
					fmt.Fprintf(os.Stdout, "\033[2J")
					Theme("header", theme)

					showWidth(filesSlice)
					showFiles(filesSlice)
					showDirs(dirs)

					bottom := "\033[" + fmt.Sprint(h) + ";0H"
					fmt.Fprintf(os.Stdout, bottom)
					Theme("footer", theme)

				}
			}
		}
	}
}
