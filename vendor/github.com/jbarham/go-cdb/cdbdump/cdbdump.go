package main

import (
	"bufio"
	"github.com/jbarham/go-cdb"
	"os"
)

func main() {
	bin, bout := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	err := cdb.Dump(bout, bin)
	bout.Flush()
	if err != nil {
		os.Exit(111)
	}
}
