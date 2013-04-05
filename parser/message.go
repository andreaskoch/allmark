// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"github.com/andreaskoch/allmark/repository"
)

func parseMessage(item *repository.Item, lines []string) *repository.Item {

	// meta data
	item, lines = ParseMetaData(item, lines)

	// content
	item.Content = getContent(lines)

	return item
}
