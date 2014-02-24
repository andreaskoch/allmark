// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

import (
	"sort"
)

type Base struct {
	Type    string `json:"type"`
	Level   int    `json:"level"`
	Route   string `json:"route"`
	BaseUrl string `json:"baseUrl"`

	Title       string `json:"title"`
	Description string `json:"description"`

	Tags []*Tag `json:"tags"`

	LanguageTag      string `json:"languageTag"`
	CreationDate     string `json:"creationdate"`
	LastModifiedDate string `json:"lastmodifieddate"`
}

type Model struct {
	Base

	Content string `json:"content"`

	TopDocs []*Model `json:"topDocs"`
	Childs  []*Base  `json:"childs"`

	ToplevelNavigation   *ToplevelNavigation   `json:"toplevelNavigation"`
	BreadcrumbNavigation *BreadcrumbNavigation `json:"breadcrumbNavigation"`

	TagCloud    *TagCloud    `json:"tagCloud"`
	Locations   []*Model     `json:"locations"`
	GeoLocation *GeoLocation `json:"geoLocation"`
}

func Error(title, content, route string) *Model {
	return &Model{
		Base: Base{
			Level:   0,
			Title:   title,
			Route:   route,
			Type:    "error",
			BaseUrl: "/",
		},
		Content: content,
	}
}

type SortModelBy func(model1, model2 *Model) bool

func (by SortModelBy) Sort(models []*Model) {
	sorter := &modelSorter{
		models: models,
		by:     by,
	}

	sort.Sort(sorter)
}

type modelSorter struct {
	models []*Model
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
