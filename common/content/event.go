// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package content

type ChangeEvent int

func (changeEvent ChangeEvent) String() string {
	switch changeEvent {
	case TypeChanged:
		return "changed"

	case TypeMoved:
		return "moved"

	default:
		return "unknown"

	}

	panic("Unreachable")
}

const (
	TypeUnknown ChangeEvent = iota
	TypeChanged
	TypeMoved
)
