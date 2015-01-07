// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"strings"
)

type Tags []Tag

func NewTags() Tags {
	return make(Tags, 0)
}

func NewTagsFromNames(names []string) (Tags, error) {
	tags := make(Tags, 0, len(names))
	errors := make([]string, 0)
	for _, name := range names {

		// skip empty values
		if strings.TrimSpace(name) == "" {
			continue
		}

		tag, err := NewTag(name)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Cannot create tag %q. Error: %s", name, err))
			continue
		}

		tags = append(tags, *tag)
	}

	return tags, fmt.Errorf("The following tags could not be created:\n%s", strings.Join(errors, "\n"))
}

func (tags Tags) Contains(otherTag Tag) bool {

	for _, tag := range tags {
		if tag.Equals(otherTag) {
			return true
		}
	}

	return false
}
