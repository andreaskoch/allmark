// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/route"
)

type ItemType int

func (itemType ItemType) String() string {
	switch itemType {
	case TypeDocument:
		return "document"

	case TypePresentation:
		return "presentation"

	case TypeMessage:
		return "message"

	case TypeLocation:
		return "location"

	case TypeRepository:
		return "repository"

	default:
		return "unknown"

	}

	panic("Unreachable")
}

const (
	TypeDocument ItemType = iota
	TypePresentation
	TypeMessage
	TypeLocation
	TypeRepository
	TypeUnknown
)

// An Item represents a single document.
type Item struct {
	isVirtual bool
	route     route.Route
	files     []*File

	Type        ItemType
	Title       string
	Description string
	Content     string

	MetaData *MetaData
}

func NewVirtualItem(route route.Route, itemType ItemType) (*Item, error) {
	return &Item{
		isVirtual: true,
		route:     route,
		Type:      itemType,
	}, nil
}

func NewItem(route route.Route, files []*File) (*Item, error) {
	return &Item{
		isVirtual: false,
		route:     route,
		files:     files,
	}, nil
}

func (item *Item) String() string {
	return fmt.Sprintf("%s", item.route)
}

func (item *Item) FolderName() string {
	return item.route.LastComponentName()
}

func (item *Item) IsVirtual() bool {
	return item.isVirtual
}

func (item *Item) Route() route.Route {
	return item.route
}

func (item *Item) Files() []*File {
	return item.files
}
