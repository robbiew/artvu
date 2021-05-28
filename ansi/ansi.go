package ansi

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	escapes "github.com/snugfox/ansi-escapes"
)

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
