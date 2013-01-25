package main

import (
	"andyk/docs/converter"
	"andyk/docs/indexer"
	"fmt"
)

func main() {
	fmt.Println(indexer.Index())
	converter.Convert()
}
