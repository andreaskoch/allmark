// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"fmt"
	"time"

	"allmark.io/modules/common/config"
	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/paths"
	"allmark.io/modules/common/route"
	"allmark.io/modules/dataaccess"
	"allmark.io/modules/model"
	"allmark.io/modules/services/converter"
	"allmark.io/modules/services/parser"
	"allmark.io/modules/web/orchestrator/index"
	"allmark.io/modules/web/orchestrator/search"
	"allmark.io/modules/web/view/viewmodel"
	"allmark.io/modules/web/webpaths"
)

type UpdateType int

const (
	UpdateTypeUnchanged UpdateType = iota
	UpdateTypeNew
	UpdateTypeModified
	UpdateTypeDeleted
)

func (updateType UpdateType) String() string {
	switch updateType {
	case UpdateTypeUnchanged:
		return "Unchanged"

	case UpdateTypeNew:
		return "new"

	case UpdateTypeModified:
		return "modified"

	case UpdateTypeDeleted:
		return "deleted"
	}

	panic("Unknown update type")
}

func NewUpdate(updateType UpdateType, route route.Route) Update {
	return Update{updateType, route}
}

type Update struct {
	updateType UpdateType
	route      route.Route
}

func (update *Update) String() string {
	return fmt.Sprintf("%s (%s)", update.route.String(), update.updateType.String())
}

func (update *Update) Route() route.Route {
	return update.route
}

func (update *Update) Type() UpdateType {
	return update.updateType
}

// cacheUpdate creates a new instance of the CacheUpdateCallback type.
func cacheUpdate(name string, updateType UpdateType, callback func(updatedRoute route.Route)) CacheUpdateCallback {
	return CacheUpdateCallback{
		name,
		updateType,
		callback,
	}
}

// CacheUpdateCallback is a wrapper model for cache update callback functions.
type CacheUpdateCallback struct {
	name       string
	updateType UpdateType
	update     func(updatedRoute route.Route)
}

// Name returns the name of the callback.
func (updateCallback *CacheUpdateCallback) Name() string {
	return updateCallback.name
}

// UpdateType returns the type of the callback.
func (updateCallback *CacheUpdateCallback) UpdateType() UpdateType {
	return updateCallback.updateType
}

// String returns a string representation of the current CacheUpdateCallback.
func (updateCallback *CacheUpdateCallback) String() string {
	return fmt.Sprintf("%s (%s)", updateCallback.Name(), updateCallback.UpdateType())
}

// Execute safely executes the callback and return any error that occured during execution.
func (updateCallback *CacheUpdateCallback) Execute(route route.Route) (err error) {
	defer func() {
		if exception := recover(); exception != nil { //catch
			err = fmt.Errorf("Error while executing callback %q. Error: %s", updateCallback.String(), exception)
		}
	}()

	updateCallback.update(route)
	return err
}

func newBaseOrchestrator(logger logger.Logger, config config.Config, repository dataaccess.Repository, parser parser.Parser, converter converter.Converter, webPathProvider webpaths.WebPathProvider) *Orchestrator {

	orchestrator := &Orchestrator{
		logger: logger,

		config:     config,
		repository: repository,
		parser:     parser,
		converter:  converter,

		webPathProvider: webPathProvider,

		updateSubscribers: make([]chan Update, 0),
		updateCallbacks:   make(map[UpdateType][]CacheUpdateCallback),
	}

	return orchestrator
}

type Orchestrator struct {
	logger logger.Logger

	config     config.Config
	repository dataaccess.Repository
	parser     parser.Parser
	converter  converter.Converter

	webPathProvider webpaths.WebPathProvider

	// caches and indizes (do not initialize!)
	fulltextIndex   *search.ItemSearch
	repositoryIndex *index.Index
	itemsByAlias    ItemCache

	// update handling
	updateCallbacks   map[UpdateType][]CacheUpdateCallback
	updateSubscribers []chan Update
}

// Get the full-page title for a given headline.
func (orchestrator *Orchestrator) GetPageTitle(headline string) string {
	rootItem := orchestrator.rootItem()
	return fmt.Sprintf("%s - %s", headline, rootItem.Title)
}

func (orchestrator *Orchestrator) Subscribe(update chan Update) {
	orchestrator.updateSubscribers = append(orchestrator.updateSubscribers, update)
}

// Update all caches
func (orchestrator *Orchestrator) UpdateCache(dataaccessLayerUpdate dataaccess.Update) {

	orchestrator.logger.Info("Received an update. Updating caches: %s", dataaccessLayerUpdate.String())

	// inform subscribers ...
	// ... about new items
	for _, newItemRoute := range dataaccessLayerUpdate.New() {

		orchestrator.logger.Info("Updating cache for route %q", newItemRoute.String())

		// execute cache update callbacks
		for _, callbackDefinition := range orchestrator.updateCallbacks[UpdateTypeNew] {
			orchestrator.logger.Debug("Executing cache update callback: %q", callbackDefinition.String())
			if err := callbackDefinition.Execute(newItemRoute); err != nil {
				orchestrator.logger.Error("%s", err.Error())
			}
		}

		// notify subscribers
		for _, subscriber := range orchestrator.updateSubscribers {
			subscriber <- NewUpdate(UpdateTypeNew, newItemRoute)
		}
	}

	// ... about modified items
	for _, modifiedItemRoute := range dataaccessLayerUpdate.Modified() {

		orchestrator.logger.Info("Updating cache for route %q", modifiedItemRoute.String())

		// execute cache update callbacks
		for _, callbackDefinition := range orchestrator.updateCallbacks[UpdateTypeModified] {
			orchestrator.logger.Debug("Executing cache update callback: %q", callbackDefinition.String())
			if err := callbackDefinition.Execute(modifiedItemRoute); err != nil {
				orchestrator.logger.Error("%s", err.Error())
			}
		}

		// notify subscribers
		for _, subscriber := range orchestrator.updateSubscribers {
			subscriber <- NewUpdate(UpdateTypeModified, modifiedItemRoute)
		}
	}

	// ... about deleted items
	for _, deletedItemRoute := range dataaccessLayerUpdate.Deleted() {

		orchestrator.logger.Info("Removing cache for route %q", deletedItemRoute.String())

		// execute cache update callbacks
		for _, callbackDefinition := range orchestrator.updateCallbacks[UpdateTypeDeleted] {
			orchestrator.logger.Debug("Executing cache update callback: %q", callbackDefinition.String())
			if err := callbackDefinition.Execute(deletedItemRoute); err != nil {
				orchestrator.logger.Error("%s", err.Error())
			}
		}

		// notify subscribers
		for _, subscriber := range orchestrator.updateSubscribers {
			subscriber <- NewUpdate(UpdateTypeDeleted, deletedItemRoute)
		}
	}

	// update the parents
	parentUpdate := getParentUpdate(dataaccessLayerUpdate)
	if parentUpdate.IsEmpty() == false {
		orchestrator.logger.Debug("Also updating the parents: %s", parentUpdate.String())
		orchestrator.UpdateCache(parentUpdate)
	}

	orchestrator.logger.Debug("Finished update (%s)", dataaccessLayerUpdate.String())
}

// registerUpdateCallback registers callbacks for new, modified and deleted items.
func (orchestrator *Orchestrator) registerUpdateCallback(name string, updateType UpdateType, callback func(updatedRoute route.Route)) {

	if orchestrator.updateCallbacks[updateType] == nil {
		orchestrator.updateCallbacks[updateType] = make([]CacheUpdateCallback, 0)
	}

	orchestrator.updateCallbacks[updateType] = append(orchestrator.updateCallbacks[updateType], cacheUpdate(name, updateType, callback))
}

func (orchestrator *Orchestrator) ItemExists(route route.Route) bool {
	_, exists := orchestrator.index().IsMatch(route)
	return exists
}

func (orchestrator *Orchestrator) absolutePather(prefix string) paths.Pather {
	return orchestrator.webPathProvider.AbsolutePather(prefix)
}

func (orchestrator *Orchestrator) itemPather() paths.Pather {
	return orchestrator.webPathProvider.ItemPather()
}

func (orchestrator *Orchestrator) tagPather() paths.Pather {
	return orchestrator.webPathProvider.TagPather()
}

func (orchestrator *Orchestrator) relativePather(baseRoute route.Route) paths.Pather {
	return orchestrator.webPathProvider.RelativePather(baseRoute)
}

func (orchestrator *Orchestrator) parseItem(item dataaccess.Item) *model.Item {
	parsedItem, err := orchestrator.parser.ParseItem(item)
	if err != nil {
		orchestrator.logger.Warn(err.Error())
		return nil
	}

	return parsedItem
}

func (orchestrator *Orchestrator) parseFile(file dataaccess.File) *model.File {
	parsedFile, err := orchestrator.parser.ParseFile(file)
	if err != nil {
		orchestrator.logger.Warn(err.Error())
		return nil
	}

	return parsedFile
}

func (orchestrator *Orchestrator) rootItem() *model.Item {
	return orchestrator.index().Root()
}

func (orchestrator *Orchestrator) getItem(route route.Route) *model.Item {

	if item, exists := orchestrator.index().IsMatch(route); exists {
		return item
	}

	return nil
}

func (orchestrator *Orchestrator) getLatestItems(parentRoute route.Route) []*model.Item {

	leafes := orchestrator.index().GetLeafes(parentRoute)

	// sort the leafes by date
	model.SortItemsBy(sortItemsByDate).Sort(leafes)

	return leafes
}

func (orchestrator *Orchestrator) index() *index.Index {

	if orchestrator.repositoryIndex != nil {
		return orchestrator.repositoryIndex
	}

	// newItem fetches the item with the given route and adds it to the index.
	updateItem := func(updatedRoute route.Route) {

		// get the updated item from the repository
		updatedRepositoryItem := orchestrator.repository.Item(updatedRoute)
		if updatedRepositoryItem == nil {
			orchestrator.logger.Warn("The item with the route %q was not found in the repository.", updatedRoute.String())
			return
		}

		// parse the item
		parsedItem := orchestrator.parseItem(updatedRepositoryItem)
		if parsedItem == nil {
			orchestrator.logger.Warn("Unable to parse item %q", updatedRepositoryItem.String())
			return
		}

		// update the index
		if orchestrator.repositoryIndex == nil {
			orchestrator.logger.Warn("Cannot add item %q, the index has not been initialized yet.", parsedItem.String())
			return
		}

		orchestrator.repositoryIndex.Add(parsedItem)
	}

	// deleteItem deletes the item with the given route from the index.
	deleteItem := func(deletedRoute route.Route) {
		orchestrator.repositoryIndex.Remove(deletedRoute)
	}

	// create a new index
	orchestrator.repositoryIndex = index.New(orchestrator.logger)

	// parse all items
	repositoryItems := orchestrator.repository.Items()
	for _, repositoryItem := range repositoryItems {
		parsedItem := orchestrator.parseItem(repositoryItem)
		if parsedItem == nil {
			orchestrator.logger.Warn("Unable to parse item %q", repositoryItem.String())
			continue
		}

		orchestrator.repositoryIndex.Add(parsedItem)
	}

	// register update callbacks
	orchestrator.registerUpdateCallback("update index", UpdateTypeNew, updateItem)
	orchestrator.registerUpdateCallback("update index", UpdateTypeModified, updateItem)
	orchestrator.registerUpdateCallback("update index", UpdateTypeDeleted, deleteItem)

	return orchestrator.repositoryIndex
}

func (orchestrator *Orchestrator) search(keywords string, maxiumNumberOfResults int) []search.Result {

	if orchestrator.fulltextIndex != nil {
		return orchestrator.fulltextIndex.Search(keywords, maxiumNumberOfResults)
	}

	// updateFulltextIndex creates a new full-text index and replaces the existing one.
	updateFulltextIndex := func(r route.Route) {
		newFullTextIndex := search.NewItemSearch(orchestrator.logger, orchestrator.getAllItems())
		orchestrator.fulltextIndex = newFullTextIndex
	}

	// initialize
	updateFulltextIndex(route.New())

	// register update callbacks
	orchestrator.registerUpdateCallback("update fulltext index", UpdateTypeNew, updateFulltextIndex)
	orchestrator.registerUpdateCallback("update fulltext index", UpdateTypeModified, updateFulltextIndex)
	orchestrator.registerUpdateCallback("update fulltext index", UpdateTypeDeleted, updateFulltextIndex)

	return orchestrator.fulltextIndex.Search(keywords, maxiumNumberOfResults)
}

func (orchestrator *Orchestrator) getAllItems() []*model.Item {

	allItems := orchestrator.index().GetAllItems()
	model.SortItemsBy(sortItemsByDate).Sort(allItems)
	return allItems

}

func (orchestrator *Orchestrator) getItems(pageSize, page int) []*model.Item {

	allItems := orchestrator.getAllItems()

	// determine the start index
	startIndex := pageSize * (page - 1)
	if startIndex >= len(allItems) {
		return []*model.Item{}
	}

	// determine the end index
	endIndex := startIndex + pageSize
	if endIndex > len(allItems) {
		endIndex = len(allItems)
	}

	return allItems[startIndex:endIndex]
}

func (orchestrator *Orchestrator) getCreationDate(itemRoute route.Route) (creationDate time.Time, found bool) {

	item := orchestrator.getItem(itemRoute)
	if item == nil {
		return time.Time{}, false
	}

	return item.MetaData.CreationDate, true
}

func (orchestrator *Orchestrator) getFile(route route.Route) *model.File {
	file, exists := orchestrator.index().IsFileMatch(route)
	if !exists {
		return nil
	}

	return file
}

func (orchestrator *Orchestrator) getParent(route route.Route) *model.Item {
	parent := orchestrator.index().GetParent(route)
	if parent == nil {
		return nil
	}

	return parent
}

func (orchestrator *Orchestrator) getPrevious(currentRoute route.Route) *model.Item {

	latestItems := orchestrator.getLatestItems(route.New())
	if len(latestItems) == 0 {
		return nil
	}

	// determine the position of the supplied route
	matchingIndex := -1
	for index, item := range latestItems {
		if item.Route().Value() == currentRoute.Value() {
			matchingIndex = index
			break
		}
	}

	// abort if the route was not found at all
	if noMatchFound := (matchingIndex == -1); noMatchFound {
		return nil
	}

	// abort if there is no next item
	nextIndex := matchingIndex + 1
	if nextIndex >= len(latestItems) {
		return nil
	}

	return latestItems[nextIndex]
}

func (orchestrator *Orchestrator) getNext(currentRoute route.Route) *model.Item {

	latestItems := orchestrator.getLatestItems(route.New())
	if len(latestItems) == 0 {
		return nil
	}

	// determine the position of the supplied route
	matchingIndex := -1
	for index, item := range latestItems {
		if item.Route().Value() == currentRoute.Value() {
			matchingIndex = index
			break
		}
	}

	// abort if the route was not found at all
	if noMatchFound := (matchingIndex == -1); noMatchFound {
		return nil
	}

	// abort if there is no previous item
	previousIndex := matchingIndex - 1
	if noPreviousItem := (previousIndex < 0); noPreviousItem {
		return nil
	}

	return latestItems[previousIndex]
}

func (orchestrator *Orchestrator) getChildren(route route.Route) []*model.Item {

	// get all children
	children := orchestrator.index().GetDirectChildren(route)

	// sort the children by date
	model.SortItemsBy(sortItemsByDate).Sort(children)

	return children
}

// getAliasMap returns the map of all items by their alias.
func (orchestrator *Orchestrator) getAliasMap() ItemCache {
	return orchestrator.itemsByAlias
}

// Get the item that has the specified alias. Returns nil if there is no matching item.
func (orchestrator *Orchestrator) getItemByAlias(alias string) *model.Item {

	// return from cache
	if orchestrator.itemsByAlias != nil {
		if item, exists := orchestrator.itemsByAlias.Get(alias); exists {
			return item
		}
		return nil
	}

	// removeItemFromAliasMap removed the item with the given route from the alias map.
	removeItemFromAliasMap := func(route route.Route) {
		// get a list of all aliases
		var aliasesToRemove []string
		for entry := range orchestrator.itemsByAlias.Iter() {
			alias := entry.Key
			item := entry.Val
			if item.Route().Equals(route) {
				aliasesToRemove = append(aliasesToRemove, alias)
			}
		}

		// remove item from map
		for _, alias := range aliasesToRemove {
			orchestrator.itemsByAlias.Remove(alias)
		}
	}

	// updateAliasMap updates the alias map for the given route.
	updateAliasMap := func(route route.Route) {

		// remove any existing aliases
		removeItemFromAliasMap(route)

		// add the new aliases
		item := orchestrator.getItem(route)
		for _, alias := range item.MetaData.Aliases {
			orchestrator.itemsByAlias.Set(alias, item)
		}
	}

	// build cache
	itemsByAlias := newItemCache()
	for _, item := range orchestrator.getAllItems() {

		for _, alias := range item.MetaData.Aliases {
			itemsByAlias.Set(alias, item)
		}
	}

	orchestrator.itemsByAlias = itemsByAlias

	// register update callbacks
	orchestrator.registerUpdateCallback("update alias map", UpdateTypeNew, updateAliasMap)
	orchestrator.registerUpdateCallback("update alias map", UpdateTypeModified, updateAliasMap)
	orchestrator.registerUpdateCallback("update alias map", UpdateTypeDeleted, removeItemFromAliasMap)

	if item, exists := orchestrator.itemsByAlias.Get(alias); exists {
		return item
	}

	return nil
}

// Get the publisher information view model.
func (orchestrator *Orchestrator) getPublisherInformation() viewmodel.Publisher {
	return viewmodel.Publisher{
		Name:  orchestrator.config.Web.Publisher.Name,
		Email: orchestrator.config.Web.Publisher.Email,
		URL:   orchestrator.config.Web.Publisher.URL,

		GooglePlusHandle: orchestrator.config.Web.Publisher.GooglePlusHandle,
		TwitterHandle:    orchestrator.config.Web.Publisher.TwitterHandle,
		FacebookHandle:   orchestrator.config.Web.Publisher.FacebookHandle,
	}
}

// Get the publisher information view model.
func (orchestrator *Orchestrator) getAuthorInformation(authorName string) viewmodel.Author {

	if authorName == "" {
		return viewmodel.Author{}
	}

	author, exists := orchestrator.config.Web.Authors[authorName]
	if !exists {
		return viewmodel.Author{
			Name: authorName,
		}
	}

	return viewmodel.Author{
		Name:  author.Name,
		Email: author.Email,
		URL:   author.URL,

		GooglePlusHandle: author.GooglePlusHandle,
		TwitterHandle:    author.TwitterHandle,
		FacebookHandle:   author.FacebookHandle,
	}
}

// Get the analytics view model.
func (orchestrator *Orchestrator) getAnalyticsSettings() viewmodel.Analytics {
	return viewmodel.Analytics{
		Enabled: orchestrator.config.Analytics.Enabled,
		GoogleAnalytics: viewmodel.GoogleAnalytics{
			Enabled:    orchestrator.config.Analytics.GoogleAnalytics.Enabled,
			TrackingID: orchestrator.config.Analytics.GoogleAnalytics.TrackingID,
		},
	}
}

// getParentUpdate returns a new Update instance that contains the parent of each item
// contained in the supplied update and marks them as "modified".
func getParentUpdate(update dataaccess.Update) dataaccess.Update {

	parentsModifiedMap := make(map[string]route.Route)

	// modified
	for _, route := range update.Modified() {
		parentRoute, exists := route.Parent()
		if !exists {
			continue
		}

		parentsModifiedMap[parentRoute.Value()] = parentRoute
	}

	// new
	for _, route := range update.New() {
		parentRoute, exists := route.Parent()
		if !exists {
			continue
		}

		parentsModifiedMap[parentRoute.Value()] = parentRoute
	}

	// deleted
	for _, route := range update.Deleted() {
		parentRoute, exists := route.Parent()
		if !exists {
			continue
		}

		parentsModifiedMap[parentRoute.Value()] = parentRoute
	}

	parentsModified := make([]route.Route, 0)
	for _, route := range parentsModifiedMap {
		parentsModified = append(parentsModified, route)
	}

	return dataaccess.NewUpdate([]route.Route{}, parentsModified, []route.Route{})
}
