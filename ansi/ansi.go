package ansi

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
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
