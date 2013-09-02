// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/view"
)

func MapTagmap(tagmap repository.TagMap) view.TagMap {

	tags := make([]*view.Tag, 0)

	for tag, items := range tagmap {

		tagModel := MapTag(tag, items)
		tags = append(tags, tagModel)
	}

	return view.TagMap{
		Tags: tags,
	}
}

func MapTag(tag repository.Tag, items repository.ItemList) *view.Tag {

	models := make([]*view.Model, 0)

	for _, item := range items {
		models = append(models, getModel(item, ""))
	}

	return &view.Tag{
		Name:          tag.Name,
		AbsoluteRoute: tag.Name,
		Description:   "not set",
		Childs:        models,
	}
}
