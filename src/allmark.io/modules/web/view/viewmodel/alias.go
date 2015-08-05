// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

import (
	"sort"
)

type AliasIndex struct {
	Model

	Aliases []Alias
}

// Alias represents an alias-index entry with a name, route and target route.
type Alias struct {
	Name        string
	Route       string
	TargetRoute string
}

// SortAliasBy can be used to sort two aliases.
type SortAliasBy func(alias1, alias2 Alias) bool

// Sort sorts the supplied list of aliases.
func (by SortAliasBy) Sort(aliases []Alias) {
	sorter := &aliasSorter{
		aliases: aliases,
		by:      by,
	}

	sort.Sort(sorter)
}

type aliasSorter struct {
	aliases []Alias
	by      SortAliasBy
}

func (sorter *aliasSorter) Len() int {
	return len(sorter.aliases)
}

func (sorter *aliasSorter) Swap(i, j int) {
	sorter.aliases[i], sorter.aliases[j] = sorter.aliases[j], sorter.aliases[i]
}

func (sorter *aliasSorter) Less(i, j int) bool {
	return sorter.by(sorter.aliases[i], sorter.aliases[j])
}
