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
