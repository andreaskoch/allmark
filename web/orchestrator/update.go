// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark/common/route"
	"github.com/andreaskoch/allmark/web/view/viewmodel"
)

type UpdateOrchestrator struct {
	*Orchestrator

	viewModelOrchestrator *ViewModelOrchestrator
}

func (orchestrator *UpdateOrchestrator) StartWatching(route route.Route) {
	orchestrator.repository.StartWatching(route)
}

func (orchestrator *UpdateOrchestrator) StopWatching(route route.Route) {
	orchestrator.repository.StopWatching(route)
}

func (orchestrator *UpdateOrchestrator) GetUpdatedModel(itemRoute route.Route) (viewModel viewmodel.Model, found bool) {
	model, found := orchestrator.viewModelOrchestrator.GetFullViewModel(itemRoute)
	if !found {
		return viewmodel.Model{}, false
	}

	return model, true
}
