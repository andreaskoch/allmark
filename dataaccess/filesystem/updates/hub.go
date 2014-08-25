// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updates

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/go-fswatch"
)

func NewHub(logger logger.Logger) *Hub {
	return &Hub{
		logger: logger,

		registry:              newRegistry(),
		triggerActionsByRoute: make(map[string]func()),
	}
}

type Hub struct {
	logger logger.Logger

	registry              *Registry
	triggerActionsByRoute map[string]func()
}

func (hub *Hub) StartWatching(route route.Route) {

	hub.logger.Debug("Watchers (Folder: %s, File: %s)", fswatch.NumberOfFolderWatchers(), fswatch.NumberOfFileWatchers())
	hub.logger.Debug("Starting callbacks for route %q", route.String())

	// get all watchers for the supplied route
	collection := hub.registry.Get(route)
	if collection == nil {
		hub.logger.Debug("There is no watcher for route %q", route.String())
		return
	}

	// start all watchers
	for _, registryEntry := range collection.Entries() {
		registryEntry.Start()
	}

	// execute the onStart trigger
	hub.executeOnStartTrigger(route)

	hub.logger.Debug("Watchers (Folder: %s, File: %s)", fswatch.NumberOfFolderWatchers(), fswatch.NumberOfFileWatchers())
}

func (hub *Hub) StopWatching(route route.Route) {

	hub.logger.Debug(fmt.Sprintf("Stopping callbacks for route %q", route.String()))

	// get all watchers for the supplied route
	collection := hub.registry.Get(route)
	if collection == nil {
		hub.logger.Debug("There is no watcher for route %q", route.String())
		return
	}

	// stop all watchers
	for _, registryEntry := range collection.Entries() {
		registryEntry.Stop()
	}
}

func (hub *Hub) Detach(route route.Route) {
	hub.logger.Debug("Detaching callbacks %q for route %q", route.String())

	collection := hub.registry.Get(route)
	if collection != nil {
		for _, entry := range collection.Entries() {
			entry.Stop()
		}
	}

	hub.registry.Remove(route)
}

func (hub *Hub) RegisterOnStartTrigger(r route.Route, action func()) {
	key := route.ToKey(r)
	hub.triggerActionsByRoute[key] = action
}

func (hub *Hub) Attach(route route.Route, callbackType string, callback func() fswatch.Watcher) {
	hub.logger.Debug("Attaching callback %q for route %q", callbackType, route.String())

	entry := newRegistryEntry(route, callbackType, callback)
	hub.registry.Add(entry)
}

func (hub *Hub) watcherExists(route route.Route, callbackType string) bool {
	collection := hub.registry.Get(route)
	if collection == nil {
		return false
	}

	key := routeAndCallbackTypeToKey(route, callbackType)
	entry := collection.Get(key)
	if entry == nil {
		return false
	}

	return true
}

func (hub *Hub) executeOnStartTrigger(r route.Route) {
	key := route.ToKey(r)
	if trigger, exists := hub.triggerActionsByRoute[key]; exists {
		go trigger()
	}
}
