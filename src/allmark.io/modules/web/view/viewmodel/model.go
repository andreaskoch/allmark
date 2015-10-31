// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

import (
	"sort"
)

type Model struct {
	Base

	Content  string `json:"content"`
	Markdown string `json:"markdown"`

	Publisher Publisher `json:"publisher"`
	Author    Author    `json:"author"`

	Children []Base `json:"children"`

	ToplevelNavigation   ToplevelNavigation   `json:"toplevelNavigation"`
	BreadcrumbNavigation BreadcrumbNavigation `json:"breadcrumbNavigation"`
	ItemNavigation       ItemNavigation       `json:"itemNavigation"`

	Tags     []Tag    `json:"tags"`
	TagCloud TagCloud `json:"tagCloud"`

	Files  []File  `json:"files"`
	Images []Image `json:"images"`

	GeoLocation GeoLocation `json:"geoLocation"`

	Analytics Analytics `json:"-"`

	Hash string `json:"hash"`

	IsRepositoryItem bool
}

func Error(title, content, route string) Model {
	return Model{
		Base: Base{
			Level:   0,
			Title:   title,
			Route:   route,
			Type:    "error",
			BaseURL: "/",
		},
		Content: content,
	}
}

type SortModelBy func(model1, model2 Model) bool

func (by SortModelBy) Sort(models []Model) {
	sorter := &modelSorter{
		models: models,
		by:     by,
	}

	sort.Sort(sorter)
}

type modelSorter struct {
	models []Model
	by     SortModelBy
}

func (sorter *modelSorter) Len() int {
	return len(sorter.models)
}

func (sorter *modelSorter) Swap(i, j int) {
	sorter.models[i], sorter.models[j] = sorter.models[j], sorter.models[i]
}

func (sorter *modelSorter) Less(i, j int) bool {
	return sorter.by(sorter.models[i], sorter.models[j])
}
