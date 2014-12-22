// Copyright 2014 Andreas Koch. All rights reserved.
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
		item := orchestrator.getLocationByName(location.Name())
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

func (orchestrator *LocationOrchestrator) getLocationByName(locationName string) *model.Item {

	item := orchestrator.getItemByAlias(locationName)
	if item != nil && item.Type == model.TypeLocation {
		return item
	}

	return nil
}

// sort tags by name
func locationModelsByFromNorthToSouth(model1, model2 *viewmodel.Model) bool {
	return model1.GeoLocation.Latitude > model2.GeoLocation.Latitude
}
