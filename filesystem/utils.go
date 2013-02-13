package filesystem

import (
	"bufio"
	"os"
)

// Readln returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned iff there is an error with the
// buffered reader.
func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)

	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}

	return string(ln), err
}

// Get all lines of a given file
func GetLines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	lines := make([]string, 0, 10)
	r := bufio.NewReader(f)
	line, err := Readln(r)
	for err == nil {
		lines = append(lines, line)
		line, err = Readln(r)
	}

	return lines, nil
}
