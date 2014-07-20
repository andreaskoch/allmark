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

func routeToKey(route route.Route) string {
	return route.Value()
}

func routeAndCallbackTypeToKey(route route.Route, callbackType string) string {
	return fmt.Sprintf("%s - %s", routeToKey(route), callbackType)
}

// Create a new registry entry from the given route, callback type and callback
func newRegistryEntry(route route.Route, callbackType string, callback func() fswatch.Watcher) *registryEntry {
	return &registryEntry{
		route:        route,
		callbackType: callbackType,

		callback: callback,
	}
}

type registryEntry struct {
	route        route.Route
	callbackType string

	callback func() fswatch.Watcher
	watcher  fswatch.Watcher
}

func (entry *registryEntry) Route() route.Route {
	return entry.route
}

func (entry *registryEntry) Type() string {
	return entry.callbackType
}

func (entry *registryEntry) Start() {

	if entry.watcher == nil {
		entry.watcher = entry.callback()
		return
	}

	if !entry.watcher.IsRunning() {
		entry.watcher.Start()
	}

}

func (entry *registryEntry) Stop() {

	if entry.watcher == nil {
		return
	}

	entry.watcher.Stop()
}

func (entry *registryEntry) Key() string {
	return routeAndCallbackTypeToKey(entry.route, entry.callbackType)
}

func (entry *registryEntry) String() string {
	return fmt.Sprintf("Registry Entry (Route: %s, Callback-Type: %s)", entry.route.Value(), entry.Type())
}

// Create a new registry entry collection
func newRegistryEntryCollection() *registryEntryCollection {
	return &registryEntryCollection{
		entriesByKey: make(map[string]*registryEntry, 0),
	}
}

// A collection of registry entries
type registryEntryCollection struct {
	entriesByKey map[string]*registryEntry
}

func (collection *registryEntryCollection) Entries() []*registryEntry {
	entries := make([]*registryEntry, 0)

	for _, value := range collection.entriesByKey {
		entries = append(entries, value)
	}

	return entries
}

func (collection *registryEntryCollection) Get(key string) *registryEntry {
	if entry, exists := collection.entriesByKey[key]; exists {
		return entry // entry was found
	}

	return nil // no entry found
}

func (collection *registryEntryCollection) Add(entry *registryEntry) bool {
	if entry := collection.Get(entry.Key()); entry != nil {
		return false // entry already exists
	}

	// add the entry to the collection
	collection.entriesByKey[entry.Key()] = entry

	return true
}

func (collection *registryEntryCollection) Remove(key string) bool {
	entry := collection.Get(key)
	if entry == nil {
		return false // there is no entry
	}

	// remove the entry
	delete(collection.entriesByKey, entry.Key())
	return true
}

func newRegistry() *Registry {
	return &Registry{
		entriesByRoute: make(map[string]*registryEntryCollection, 0),
	}
}

type Registry struct {
	entriesByRoute map[string]*registryEntryCollection
}

func (registry *Registry) Get(route route.Route) *registryEntryCollection {
	if collection, exists := registry.entriesByRoute[routeToKey(route)]; exists {
		return collection
	}

	return nil
}

func (registry *Registry) Add(entry *registryEntry) bool {
	collection := registry.Get(entry.Route())
	if collection == nil {
		collection = newRegistryEntryCollection() // create a new collection if the is none for the given route
		registry.entriesByRoute[routeToKey(entry.Route())] = collection
	}

	return collection.Add(entry)
}

func (registry *Registry) Remove(route route.Route) bool {

	// check if there is a collection for the given route
	collection := registry.Get(route)
	if collection == nil {
		return false // there was no collection for the supplied route
	}

	// remove the entry
	key := routeToKey(route)
	delete(registry.entriesByRoute, key)

	return true
}

func newUpdateHub(logger logger.Logger) *UpdateHub {
	return &UpdateHub{
		logger: logger,

		registry: newRegistry(),
	}
}

type UpdateHub struct {
	logger logger.Logger

	registry *Registry
}

func (hub *UpdateHub) StartWatching(route route.Route) {

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

	hub.logger.Debug("Watchers (Folder: %s, File: %s)", fswatch.NumberOfFolderWatchers(), fswatch.NumberOfFileWatchers())
}

func (hub *UpdateHub) StopWatching(route route.Route) {

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

func (hub *UpdateHub) Detach(route route.Route) {
	hub.logger.Debug("Detaching callbacks %q for route %q", route.String())

	collection := hub.registry.Get(route)
	if collection != nil {
		for _, entry := range collection.Entries() {
			entry.Stop()
		}
	}

	hub.registry.Remove(route)
}

func (hub *UpdateHub) Attach(route route.Route, callbackType string, callback func() fswatch.Watcher) {
	hub.logger.Debug("Attaching callback %q for route %q", callbackType, route.String())

	entry := newRegistryEntry(route, callbackType, callback)
	hub.registry.Add(entry)
}

func (hub *UpdateHub) watcherExists(route route.Route, callbackType string) bool {
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
