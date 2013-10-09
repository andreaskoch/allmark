// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

type LocationMap map[Location]ItemList

func NewLocationMap() LocationMap {
	return make(LocationMap)
}

func (locationmap LocationMap) Register(item *Item, previousLocationList Locations) {

	// add new location
	for _, location := range item.MetaData.Locations {

		if itemlist, exists := locationmap[location]; exists {

			// add the item to the item list for this location
			locationmap[location] = itemlist.Add(item)

		} else {

			// create a new item list
			locationmap[location] = NewItemList(item)
		}

	}

	// remove old location
	for _, location := range previousLocationList {

		// check if the old location is still in the new location list
		if item.MetaData.Locations.Contains(location) {
			continue // the location is still there
		}

		// the location has been removed from the item's location list
		if itemlist, exists := locationmap[location]; exists {

			// remove the item from the item list for this location
			locationmap[location] = itemlist.Remove(item)

		}
	}

}

func (locationmap LocationMap) Remove(item *Item) {

	for location, itemlist := range locationmap {

		locationmap[location] = itemlist.Remove(item)

		// remove location if item list is empty
		if itemlist.IsEmpty() {
			delete(locationmap, location)
		}
	}

}
