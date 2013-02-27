package server

import (
	"andyk/docs/indexer"
)

func Serve(repositoryPaths []string) {
	for _, repositoryPath := range repositoryPaths {
		indexer.GetIndex(repositoryPath)

	}
}
