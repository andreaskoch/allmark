// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package files

import (
	"github.com/andreaskoch/allmark2/common/route"
)

type FileSystemEntry struct {
	name   string
	parent *FileSystemEntry
	Childs []*FileSystemEntry
}

func NewRootFilesystemEntry(name string) *FileSystemEntry {
	return &FileSystemEntry{
		name:   name,
		Childs: make([]*FileSystemEntry, 0),
	}
}

func NewFilesystemEntry(parent *FileSystemEntry, name string) *FileSystemEntry {
	return &FileSystemEntry{
		name:   name,
		parent: parent,
		Childs: make([]*FileSystemEntry, 0),
	}
}

func (fsEntry *FileSystemEntry) Path() string {
	path := fsEntry.name

	if fsEntry.Parent() != nil {
		path = fsEntry.Parent().Path() + "/" + path
	}

	return path
}

func (fsEntry *FileSystemEntry) IsDirectory() bool {
	return len(fsEntry.Childs) > 0
}

func (fsEntry *FileSystemEntry) Parent() *FileSystemEntry {
	return fsEntry.parent
}

func (fsEntry *FileSystemEntry) Name() string {
	return route.DecodeUrl(fsEntry.name)
}

func (fsEntry *FileSystemEntry) GetChild(name string) *FileSystemEntry {
	for _, entry := range fsEntry.Childs {
		if entry.name == name {
			return entry
		}
	}

	return nil
}
