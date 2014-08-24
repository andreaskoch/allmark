// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
)

type UpdateOrchestrator struct {
	*Orchestrator

	viewModelOrchestrator ViewModelOrchestrator
}

func (orchestrator *UpdateOrchestrator) StartWatching(route route.Route) {
	orchestrator.repository.StartWatching(route)
}

func (orchestrator *UpdateOrchestrator) StopWatching(route route.Route) {
	orchestrator.repository.StopWatching(route)
}

func (orchestrator *UpdateOrchestrator) OnUpdate(callback func(route.Route)) {
	orchestrator.logger.Info("Assigning a new update-callback to the repository.")
	orchestrator.repository.OnUpdate(callback)
}

func (orchestrator *UpdateOrchestrator) GetUpdatedModel(itemRoute route.Route) (viewModel viewmodel.Model, found bool) {
	return orchestrator.viewModelOrchestrator.GetViewModel(itemRoute)
}
