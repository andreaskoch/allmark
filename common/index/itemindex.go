// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/model"
	"github.com/bradleypeabody/fulltext"
)

func CreateItemIndex(logger logger.Logger) *ItemIndex {
	return &ItemIndex{
		logger: logger,
		items:  make(map[route.Route]*model.Item),
	}
}

type ItemIndex struct {
	logger logger.Logger
	items  map[route.Route]*model.Item
}

func (index *ItemIndex) IsMatch(route route.Route) (item *model.Item, isMatch bool) {

	// check for a direct match
	if item, isMatch = index.items[route]; isMatch {
		return item, isMatch
	}

	// no match
	return nil, false
}

func (index *ItemIndex) IsFileMatch(route route.Route) (*model.File, bool) {

	var parent *model.Item
	parentRoute := &route
	for parentRoute != nil && parentRoute.Level() > 0 {

		parent, _ = index.IsMatch(*parentRoute)
		if parent == nil || parent.IsVirtual() {

			// next level
			parentRoute = parentRoute.Parent()
			continue

		}

		// found a non-virtual parent
		break

	}

	// abort if there is no non-virtual parent
	if parent == nil || parent.IsVirtual() {
		return nil, false
	}

	// check if the parent has a file with the supplied route
	if file := parent.GetFile(route); file != nil {
		return file, true
	}

	// file not found
	return nil, false
}

func (index *ItemIndex) GetParent(childRoute *route.Route) *model.Item {

	if childRoute == nil {
		return nil
	}

	// abort if the supplied route is already a root
	if childRoute.Level() == 0 {
		return nil
	}

	// get the parent route
	parentRoute := childRoute.Parent()
	if parentRoute == nil {
		return nil
	}

	item, isMatch := index.IsMatch(*parentRoute)
	if !isMatch {
		return nil
	}

	return item
}

func (index *ItemIndex) Root() *model.Item {
	root := route.New()
	return index.items[*root]
}

func (index *ItemIndex) GetAllChilds(route *route.Route) []*model.Item {
	return index.getChilds(route, true)
}

func (index *ItemIndex) GetChilds(route *route.Route) []*model.Item {
	return index.getChilds(route, false)
}

func (index *ItemIndex) getChilds(route *route.Route, recurse bool) []*model.Item {

	routeLevel := route.Level()
	nextLevel := routeLevel + 1

	// routeLevel := route.Level()
	childs := make([]*model.Item, 0)

	for childRoute, child := range index.items {

		// skip all deeper-level childs if recursion is disabled
		if !recurse && childRoute.Level() != nextLevel {
			continue
		}

		// skip all items which are not a child
		if !childRoute.IsChildOf(route) {
			continue
		}

		childs = append(childs, child)
	}

	// sort the items by ascending by route
	model.SortItemBy(sortItemsByRoute).Sort(childs)

	return childs
}

func (index *ItemIndex) Routes() []route.Route {
	routes := make([]route.Route, 0)
	for route, _ := range index.items {
		routes = append(routes, route)
	}
	return routes
}

func (index *ItemIndex) Items() []*model.Item {
	items := make([]*model.Item, 0, len(index.items))

	for _, item := range index.items {
		items = append(items, item)
	}

	// sort the items by ascending by route
	model.SortItemBy(sortItemsByRoute).Sort(items)

	return items
}

// Get the maxium level of all routes in this index (default: 0)
func (index *ItemIndex) MaxLevel() int {

	maxLevel := 0

	for _, item := range index.items {
		itemLevel := item.Route().Level()
		if itemLevel > maxLevel {
			maxLevel = itemLevel
		}
	}

	return maxLevel
}

func (index *ItemIndex) Add(item *model.Item) {
	index.logger.Debug("Adding item %q to index", item)

	// the the item to the index
	itemRoute := *item.Route()
	index.items[itemRoute] = item

	// insert virtual items if required
	index.fillGapsWithVirtualItems(itemRoute)

	// fulltext search
	idx, err := fulltext.NewIndexer("")
	if err != nil {
		panic(err)
	}
	defer idx.Close()

	// for each document you want to add, you do something like this:
	doc := fulltext.IndexDoc{
		Id:         []byte(itemRoute.Value()), // unique identifier (the path to a webpage works...)
		StoreValue: []byte(item.Title),        // bytes you want to be able to retrieve from search results
		IndexValue: []byte(item.Title),        // bytes you want to be split into words and indexed
	}
	idx.AddDoc(doc) // add it

	// when done, write out to final index
	f, err := fsutil.OpenFile("index")
	if err != nil {
		panic(err)
	}

	err = idx.FinalizeAndWrite(f)
	if err != nil {
		panic(err)
	}
}

func (index *ItemIndex) Search(keyword string) []SearchResult {

	searcher, err := fulltext.NewSearcher("index")
	if err != nil {
		panic(err)
	}

	defer searcher.Close()

	searchResult, err := searcher.SimpleSearch(keyword, 5)
	if err != nil {
		panic(err)
	}

	index.logger.Debug("%s", keyword)
	index.logger.Debug("%s", len(searchResult.Items))

	searchResults := make([]SearchResult, 0)

	for k, v := range searchResult.Items {
		index.logger.Info("----------- #:%d\n", k)
		index.logger.Info("Id: %s\n", v.Id)
		index.logger.Info("Score: %d\n", v.Score)
		index.logger.Info("StoreValue: %s\n", v.StoreValue)

		route, err := route.NewFromRequest(string(v.Id))
		if err != nil {
			index.logger.Warn("%s", err)
			continue
		}

		if item, isMatch := index.IsMatch(*route); isMatch {
			searchResults = append(searchResults, SearchResult{
				Score:      v.Score,
				StoreValue: string(v.StoreValue),
				Item:       item,
			})
		}

	}

	return searchResults
}

type SearchResult struct {
	Score      int64
	StoreValue string
	Item       *model.Item
}

func (index *ItemIndex) fillGapsWithVirtualItems(baseRoute route.Route) {

	// validate the input
	if baseRoute.Level() == 0 {
		index.logger.Debug("%q is at level 0", baseRoute)
		return
	}

	parentRoute := baseRoute.Parent()
	for parentRoute != nil && parentRoute.Level() > 0 {

		if _, exists := index.items[*parentRoute]; !exists {

			index.logger.Debug("Adding virtual item %q to index", parentRoute)

			virtualParentItem, err := newVirtualItem(*parentRoute)
			if err != nil {
				panic(err)
			}

			// add the virtual item to the index
			index.items[*parentRoute] = virtualParentItem

		}

		// move up
		parentRoute = parentRoute.Parent()

	}
}

func newVirtualItem(route route.Route) (*model.Item, error) {

	// create a virtual item
	item, err := model.NewVirtualItem(&route)
	if err != nil {
		return nil, err
	}

	// set the item title
	item.Title = route.FolderName()

	return item, nil
}

// sort the items by name
func sortItemsByRoute(item1, item2 *model.Item) bool {

	// ascending by route
	return item1.Route().Value() < item2.Route().Value()

}
