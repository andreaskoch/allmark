// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/dataaccess"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/converter"
	"github.com/andreaskoch/allmark2/services/parser"
	"github.com/andreaskoch/allmark2/web/webpaths"
	"strings"
	"time"
)

func newBaseOrchestrator(logger logger.Logger, repository dataaccess.Repository, parser parser.Parser, converter converter.Converter, webPathProvider webpaths.WebPathProvider) *Orchestrator {
	return &Orchestrator{
		logger: logger,

		repository: repository,
		parser:     parser,
		converter:  converter,

		webPathProvider: webPathProvider,
	}
}

type Orchestrator struct {
	logger logger.Logger

	repository dataaccess.Repository
	parser     parser.Parser
	converter  converter.Converter

	webPathProvider webpaths.WebPathProvider

	itemsByAlias map[string]*model.Item
}

func (orchestrator *Orchestrator) ItemExists(route route.Route) bool {
	_, exists := orchestrator.repository.Item(route)
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

func (orchestrator *Orchestrator) parseItem(item *dataaccess.Item) *model.Item {
	parsedItem, err := orchestrator.parser.ParseItem(item)
	if err != nil {
		orchestrator.logger.Warn(err.Error())
		return nil
	}

	return parsedItem
}

func (orchestrator *Orchestrator) parseFile(file *dataaccess.File) *model.File {
	parsedFile, err := orchestrator.parser.ParseFile(file)
	if err != nil {
		orchestrator.logger.Warn(err.Error())
		return nil
	}

	return parsedFile
}

func (orchestrator *Orchestrator) rootItem() *model.Item {
	return orchestrator.parseItem(orchestrator.repository.Root())
}

func (orchestrator *Orchestrator) getItem(route route.Route) *model.Item {
	item, exists := orchestrator.repository.Item(route)
	if !exists {
		return nil
	}

	return orchestrator.parseItem(item)
}

func (orchestrator *Orchestrator) getAllItems() []*model.Item {

	allItems := make([]*model.Item, 0)

	for _, repositoryItem := range orchestrator.repository.Items() {
		item := orchestrator.parseItem(repositoryItem)
		if item == nil {
			continue
		}

		allItems = append(allItems, item)
	}

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

func (orchestrator *ViewModelOrchestrator) getCreationDate(itemRoute route.Route) (creationDate time.Time, found bool) {

	item := orchestrator.getItem(itemRoute)
	if item == nil {
		return time.Time{}, false
	}

	return item.MetaData.CreationDate, true
}

func (orchestrator *Orchestrator) getFile(route route.Route) *model.File {
	file, exists := orchestrator.repository.File(route)
	if !exists {
		return nil
	}

	return orchestrator.parseFile(file)
}

func (orchestrator *Orchestrator) getParent(route route.Route) *model.Item {
	parent := orchestrator.repository.Parent(route)
	if parent == nil {
		return nil
	}

	return orchestrator.parseItem(parent)
}

func (orchestrator *Orchestrator) getChilds(route route.Route) (childs []*model.Item) {

	childs = make([]*model.Item, 0)

	for _, child := range orchestrator.repository.Childs(route) {
		parsed := orchestrator.parseItem(child)
		if parsed == nil {
			continue
		}

		childs = append(childs, parsed)
	}

	return childs
}

// Get the item that has the specified alias. Returns nil if there is no matching item.
func (orchestrator *Orchestrator) getItemByAlias(alias string) *model.Item {

	alias = strings.TrimSpace(strings.ToLower(alias))

	if orchestrator.itemsByAlias == nil {

		orchestrator.logger.Info("Initializing alias list")
		itemsByAlias := make(map[string]*model.Item)

		for _, repositoryItem := range orchestrator.repository.Items() {

			item := orchestrator.parseItem(repositoryItem)
			if item == nil {
				orchestrator.logger.Warn("Cannot parse repository item %q.", repositoryItem.String())
				continue
			}

			// continue items without an alias
			if item.MetaData.Alias == "" {
				continue
			}

			itemAlias := strings.TrimSpace(strings.ToLower(item.MetaData.Alias))
			itemsByAlias[itemAlias] = item
		}

		// refresh control
		go func() {
			for {
				select {
				case <-orchestrator.repository.AfterReindex():
					// reset the list
					orchestrator.logger.Info("Resetting the alias list")
					orchestrator.itemsByAlias = nil
				}
			}
		}()

		orchestrator.itemsByAlias = itemsByAlias
	}

	return orchestrator.itemsByAlias[alias]
}

// sort the models by date and name
func sortItemsByDate(model1, model2 *model.Item) bool {

	return model1.MetaData.CreationDate.Before(model2.MetaData.CreationDate)

}
