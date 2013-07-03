// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"fmt"
	"github.com/andreaskoch/allmark/config"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/util"
	"path/filepath"
)

type Indexer struct {
	RootIsReady chan *Item
	New         chan *Item
	Deleted     chan *Item

	indexPath            string
	relativePathProvider *path.Provider
	config               *config.Config
}

func New(indexPath string, config *config.Config, useTempDir bool) (*Indexer, error) {

	// check if path exists
	if !util.PathExists(indexPath) {
		return nil, fmt.Errorf("The path %q does not exist.", indexPath)
	}

	// check if the path is a directory
	if isDirectory, _ := util.IsDirectory(indexPath); !isDirectory {
		indexPath = filepath.Dir(indexPath)
	}

	if isReservedDirectory(indexPath) {
		return nil, fmt.Errorf("The path %q is using a reserved name and cannot be a root.", indexPath)
	}

	// create a new indexer
	indexer := &Indexer{
		RootIsReady: make(chan *Item),
		New:         make(chan *Item),
		Deleted:     make(chan *Item),

		indexPath:            indexPath,
		relativePathProvider: path.NewProvider(indexPath, useTempDir),
		config:               config,
	}

	return indexer, nil
}

func (indexer *Indexer) RelativePathProvider() *path.Provider {
	return indexer.relativePathProvider
}

func (indexer *Indexer) Execute() {

	// create a new item
	rootItem, err := newItem(indexer.RelativePathProvider(), nil, indexer.indexPath, 0, indexer.New, indexer.Deleted)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case <-rootItem.ChildsReady:
				go func() {
					indexer.New <- rootItem
					indexer.RootIsReady <- rootItem
				}()
			}
		}
	}()
}
