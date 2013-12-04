// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package message

import (
	"fmt"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/parser/metadata"
	"strings"
	"time"
)

func Parse(item *model.Item, lastModifiedDate time.Time, lines []string) (parseError error) {

	// content
	contentStartIndex := 0
	contentEndIndex := len(lines)

	if metaDataStartIndex, err := metadata.GetMetaDataPosition(lines); err == nil {
		contentEndIndex = metaDataStartIndex
	}
	contentLines := lines[contentStartIndex:contentEndIndex]

	// message
	message := strings.TrimSpace(strings.Join(contentLines, "\n"))
	item.Description = message
	item.Content = message

	// meta data
	if err := metadata.Parse(item, lastModifiedDate, lines); err != nil {
		return fmt.Errorf("Unable to parse the meta data of item %q. Error: %s", item, err)
	}

	// title
	item.Title = item.MetaData.CreationDate.Format(time.RFC850)

	return
}
