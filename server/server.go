package server

import (
	"andyk/docs/indexer"
	"andyk/docs/renderer"
	"fmt"
)

func Serve(repositoryPaths []string) {

	itemPaths := make([]string, 0, 0)
	for _, repositoryPath := range repositoryPaths {

		// create an index
		index := indexer.GetIndex(repositoryPath)

		// render all index items
		renderer.RenderIndex(index)

		itemPaths = append(itemPaths, index.GetRelativeItemPaths()...)
	}

	fmt.Printf("%#v", itemPaths)
}
