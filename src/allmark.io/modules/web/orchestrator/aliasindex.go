// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"allmark.io/modules/web/view/viewmodel"
	"fmt"
)

// AliasIndexOrchestrator provides alias index entries.
type AliasIndexOrchestrator struct {
	*Orchestrator
}

// GetIndexEntries returns a list of all alias index entry models.
func (orchestrator *AliasIndexOrchestrator) GetIndexEntries(hostname, prefix string) []viewmodel.Alias {

	itemPathProvider := orchestrator.absolutePather("/")
	aliasPathProvider := orchestrator.absolutePather("/" + prefix)

	var aliasIndexEntries []viewmodel.Alias

	for entry := range orchestrator.getAliasMap().Iter() {
		alias := entry.Key
		item := entry.Val

		aliasIndexEntries = append(aliasIndexEntries, viewmodel.Alias{
			Name:        fmt.Sprintf("%s%s", prefix, alias),
			Route:       aliasPathProvider.Path(alias),
			TargetRoute: itemPathProvider.Path(item.Route().Value()),
		})

	}

	// sort aliases by name
	viewmodel.SortAliasBy(aliasByName).Sort(aliasIndexEntries)

	return aliasIndexEntries
}

// aliasByName returns true if alias 1 comes before alias 2; otherwise false.
func aliasByName(alias1, alias2 viewmodel.Alias) bool {
	return alias1.Name < alias2.Name
}
