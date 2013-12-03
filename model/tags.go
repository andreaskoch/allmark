// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
)

type Tags []Tag

func NewTags() Tags {
	return make(Tags, 0)
}

func NewTagsFromNames(names []string) Tags {
	tags := make(Tags, 0, len(names))

	for _, name := range names {

		tag, err := NewTag(name)
		if err != nil {
			fmt.Printf("Skipping tag %q. Error: %s\n", name, err)
			continue
		}

		tags = append(tags, *tag)
	}

	return tags
}

func (tags Tags) Contains(otherTag Tag) bool {

	for _, tag := range tags {
		if tag.Equals(otherTag) {
			return true
		}
	}

	return false
}
