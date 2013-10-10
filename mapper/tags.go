// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/view"
)

var (
	// sort tags by name
	tagsByName = func(tag1, tag2 *view.Tag) bool {
		return tag1.Name < tag2.Name
	}
)

func MapTagmap(tagmap repository.TagMap, relativePath func(item *repository.Item) string, absolutePath func(item *repository.Item) string, content func(item *repository.Item) string) view.TagMap {

	tags := make([]*view.Tag, 0)

	for tag, items := range tagmap {

		tagModel := MapTag(tag, items, relativePath, absolutePath, content)
		tags = append(tags, tagModel)
	}

	// sort tags by name
	view.SortTagBy(tagsByName).Sort(tags)

	return view.TagMap{
		Tags: tags,
	}
}

func MapTag(tag repository.Tag, items repository.ItemList, relativePath func(item *repository.Item) string, absolutePath func(item *repository.Item) string, content func(item *repository.Item) string) *view.Tag {

	models := make([]*view.Model, 0)

	for _, item := range items {
		models = append(models, Map(item, relativePath, absolutePath, content))
	}

	return &view.Tag{
		Name:          tag.Name(),
		AbsoluteRoute: tagPathResolver(&tag),
		Description:   "",
		Childs:        models,
	}
}
