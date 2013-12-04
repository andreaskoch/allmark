// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dataaccess

type RepositoryEvent struct {
	Item  *Item
	Error error
}

func NewEvent(item *Item, err error) *RepositoryEvent {
	return &RepositoryEvent{
		Item:  item,
		Error: err,
	}
}

type Repository interface {
	GetItems() (itemEvents chan *RepositoryEvent, done chan bool)
	Id() string
	Path() string
}
