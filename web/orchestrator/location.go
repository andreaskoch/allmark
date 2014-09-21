// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
)

type LocationOrchestrator struct {
	*Orchestrator

	locationsByAlias map[string]*model.Item
}

func (orchestrator *LocationOrchestrator) GetLocations(locations model.Locations, conversion func(*model.Item) viewmodel.Model) []*viewmodel.Model {
	locationModels := make([]*viewmodel.Model, 0)

	for _, location := range locations {
		item := orchestrator.getItemFromLocationName(location.Name())
		if item == nil {
			orchestrator.logger.Warn("Location %q was not found.", location.Name())
			continue
		}

		viewmodel := conversion(item)
		locationModels = append(locationModels, &viewmodel)
	}

	// sort locations from north to south
	viewmodel.SortModelBy(locationModelsByFromNorthToSouth).Sort(locationModels)

	return locationModels
}

func (orchestrator *LocationOrchestrator) getItemFromLocationName(locationName string) *model.Item {
	return orchestrator.getLocationByName(locationName)
}

func (orchestrator *LocationOrchestrator) getLocationByName(locationName string) *model.Item {

	if orchestrator.locationsByAlias == nil {

		locationsByAlias := make(map[string]*model.Item)

		for _, repositoryItem := range orchestrator.repository.Items() {

			item := orchestrator.parseItem(repositoryItem)
			if item == nil {
				orchestrator.logger.Warn("Cannot parse repository item %q.", repositoryItem.String())
				continue
			}

			// skip non-location items
			if item.Type != model.TypeLocation {
				continue
			}

			// skip items without meta data
			if item.MetaData == nil {
				continue
			}

			// continue items without an alias
			if item.MetaData.Alias == "" {
				continue
			}

			locationsByAlias[item.MetaData.Alias] = item
		}

		orchestrator.locationsByAlias = locationsByAlias
	}

	return orchestrator.locationsByAlias[locationName]
}

// sort tags by name
func locationModelsByFromNorthToSouth(model1, model2 *viewmodel.Model) bool {
	return model1.GeoLocation.Latitude > model2.GeoLocation.Latitude
}
