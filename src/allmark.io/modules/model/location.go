// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"strings"
)

type Location struct {
	name string
}

func NewLocation(name string) (*Location, error) {

	normalized := normalizeLocationName(name)
	if normalized == "" {
		return nil, fmt.Errorf("Cannot create a location from an empty string")
	}

	return &Location{
		name: normalized,
	}, nil
}

func (location *Location) String() string {
	return location.name
}

func (location *Location) Name() string {
	return location.name
}

func (location *Location) Equals(otherLocation Location) bool {
	return location.Name() == otherLocation.Name()
}

func normalizeLocationName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}
