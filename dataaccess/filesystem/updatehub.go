// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/go-fswatch"
)

type updateHubCallbacks map[string]func() fswatch.Watcher
type updateHubWatchers map[string]fswatch.Watcher

type watcherRegistry map[string]updateHubWatchers
type callbackRegistry map[string]updateHubCallbacks

func newUpdateHub(logger logger.Logger) *UpdateHub {
	return &UpdateHub{
		logger: logger,

		callbacks: make(callbackRegistry),
		watchers:  make(watcherRegistry),
	}
}

type UpdateHub struct {
	logger logger.Logger

	callbacks callbackRegistry
	watchers  watcherRegistry
}

func (hub *UpdateHub) StartWatching(route route.Route) {

	hub.logger.Debug("# folder watchers: %v", fswatch.NumberOfFolderWatchers())

	hub.logger.Debug(fmt.Sprintf("Starting callbacks for route %q", route.String()))

	for callbackType, callback := range hub.callbacks[routeToKey(route)] {
		hub.logger.Debug(fmt.Sprintf("Starting callback %q for route %q", callbackType, route.String()))

		// execute the callback
		watcher := callback()

		if watchers, exists := hub.watchers[routeToKey(route)]; !exists {
			watchers := make(updateHubWatchers)
			watchers[callbackType] = watcher
			hub.watchers[routeToKey(route)] = watchers
		} else {
			watchers[callbackType] = watcher
			hub.watchers[routeToKey(route)] = watchers
		}
	}
}

func (hub *UpdateHub) StopWatching(route route.Route) {

	hub.logger.Debug(fmt.Sprintf("Stopping callbacks for route %q", route.String()))

	watchers, exists := hub.watchers[routeToKey(route)]
	if !exists {
		hub.logger.Debug("There is no running watcher for route %q", route.String())
		return
	}

	for callbackType, watcher := range watchers {
		hub.logger.Debug("Stopping watcher %q for route %q", callbackType, route.String())
		if watcher != nil {
			watcher.Stop()
		}
	}

}

func (hub *UpdateHub) Detach(route route.Route) {
	hub.logger.Debug("Detaching callbacks %q for route %q", route.String())
	hub.StopWatching(route)
	delete(hub.watchers, route.Value())
}

func (hub *UpdateHub) Attach(route route.Route, callbackType string, callback func() fswatch.Watcher) {
	hub.logger.Debug("Attaching callback %q for route %q", callbackType, route.String())

	if callbacks, exists := hub.callbacks[routeToKey(route)]; !exists {

		// create a new callback map
		callbacks := make(updateHubCallbacks)

		// attach the callback
		callbacks[callbackType] = callback

		hub.callbacks[routeToKey(route)] = callbacks
	} else {

		// attach the callback
		callbacks[callbackType] = callback

	}
}

func routeToKey(route route.Route) string {
	return route.Value()
}
