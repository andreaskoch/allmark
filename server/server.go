package main

import (
	"andyk/docs/indexer"
	"andyk/docs/converter"
)

func main() {
	indexer.Index()
	converter.Convert()
}