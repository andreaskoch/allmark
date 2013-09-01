// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"fmt"
)

type Tags []Tag

func NewTags(names []string) Tags {
	tags := Tags{}

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
