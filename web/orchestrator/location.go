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
}

func (orchestrator *LocationOrchestrator) GetLocations(locations model.Locations) []*viewmodel.Model {
	locationModels := make([]*viewmodel.Model, 0)

	for _, location := range locations {
		item := orchestrator.getItemFromLocationName(location.Name())
		if item == nil {
			orchestrator.logger.Warn("Location %q was not found.", location.Name())
			continue
		}

		locationModels = append(locationModels, nil)
	}

	// sort locations from north to south
	viewmodel.SortModelBy(locationModelsByFromNorthToSouth).Sort(locationModels)

	return locationModels
}

func (Orchestrator *LocationOrchestrator) getItemFromLocationName(locationName string) *model.Item {
	return nil
}

// sort tags by name
func locationModelsByFromNorthToSouth(model1, model2 *viewmodel.Model) bool {
	return model1.GeoLocation.Latitude > model2.GeoLocation.Latitude
}
