// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"fmt"
)

type Locations []Location

func NewLocations() Locations {
	return make(Locations, 0)
}

func NewLocationsFromNames(names []string) Locations {
	locations := make(Locations, 0, len(names))

	for _, name := range names {

		location, err := NewLocation(name)
		if err != nil {
			fmt.Printf("Skipping location %q. Error: %s\n", name, err)
			continue
		}

		locations = append(locations, *location)
	}

	return locations
}

func (locations Locations) Contains(otherLocation Location) bool {

	for _, location := range locations {
		if location.Equals(otherLocation) {
			return true
		}
	}

	return false
}
