// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dataaccess

import (
	"github.com/elWyatt/allmark/common/content"
	"github.com/elWyatt/allmark/common/route"
)

type ItemType int

func (itemType ItemType) String() string {
	switch itemType {

	case TypePhysical:
		return "physical"

	case TypeVirtual:
		return "virtual"

	case TypeFileCollection:
		return "filecollection"

	default:
		return "unknown"

	}

	panic("Unreachable")
}

const (
	TypePhysical ItemType = iota
	TypeVirtual
	TypeFileCollection
)

type ItemState int

const (
	ItemStateStable ItemState = iota
	ItemStateNew
	ItemStateModified
	ItemStateDeleted
)

// An Item represents a single document in a repository.
type Item interface {
	content.ContentProviderInterface

	String() string
	Id() string
	Type() ItemType
	CanHaveChildren() bool
	Route() route.Route
	Files() []File
	LastHash() string
}
