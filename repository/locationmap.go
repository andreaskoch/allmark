// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

type LocationResolver func(locationName string) ItemList

type LocationMap map[string]ItemList

func NewLocationMap() LocationMap {
	return make(LocationMap)
}

func (locationmap LocationMap) Register(item *Item, previousLocationList Locations) {

	// add new location
	for _, location := range item.MetaData.Locations {

		key := location.Name()
		if itemlist, exists := locationmap[key]; exists {

			// add the item to the item list for this location
			locationmap[key] = itemlist.Add(item)

		} else {

			// create a new item list
			locationmap[key] = NewItemList(item)
		}

	}

	// remove old location
	for _, location := range previousLocationList {

		key := location.Name()

		// check if the old location is still in the new location list
		if item.MetaData.Locations.Contains(location) {
			continue // the location is still there
		}

		// the location has been removed from the item's location list
		if itemlist, exists := locationmap[key]; exists {

			// remove the item from the item list for this location
			locationmap[key] = itemlist.Remove(item)

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

func (locationmap LocationMap) Lookup(locationName string) ItemList {

	if location, err := NewLocation(locationName); err == nil {
		key := location.Name()
		return locationmap[key]
	}

	return nil
}
