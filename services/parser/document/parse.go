// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package document

import (
	"fmt"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/parser/metadata"
	"github.com/andreaskoch/allmark2/services/parser/pattern"
	"strings"
	"time"
)

func Parse(item *model.Item, lastModifiedDate time.Time, lines []string) (warning, err error) {

	// title
	titleLineNumber := len(lines)
	for lineNumber, line := range lines {

		// ignore empty lines
		if pattern.IsEmpty(line) {
			continue
		}

		// search for the title
		if isTitle, title := pattern.IsTitle(line); isTitle {
			// capture the line number of the title
			titleLineNumber = lineNumber

			// save the title to the item
			item.Title = strings.TrimSpace(title)

		} else {

			// assign a fallback title
			item.Title = item.FolderName()

			// reuse this line for the description or content
			titleLineNumber = -1

		}

		break
	}

	// abort if there are no more lines
	if len(lines) < titleLineNumber+1 {
		return nil, nil
	}

	// description
	descriptionLineNumber := titleLineNumber + 1
	for lineNumber, line := range lines[(titleLineNumber + 1):] {

		// ignore empty lines
		if pattern.IsEmpty(line) {
			continue
		}

		// search for the description
		if isDescription, description := pattern.IsDescription(line); isDescription {

			// capture the line number of the description
			descriptionLineNumber = lineNumber

			// save the description to the item
			item.Description = strings.TrimSpace(description)

		} else {

			// reuse this line for the content. 2 because a description is supposed to be followed by an empty line
			descriptionLineNumber = lineNumber - 2

		}

		break
	}

	// abort if there are no more lines
	if len(lines) < descriptionLineNumber+2 {
		return // there are no more lines, but having no content is ok
	}

	// content
	contentStartIndex := (descriptionLineNumber + 2)
	contentEndIndex := len(lines)

	if metaDataStartIndex, err := metadata.GetMetaDataPosition(lines); err == nil {
		contentEndIndex = metaDataStartIndex
	}

	contentLines := lines[contentStartIndex:contentEndIndex]
	item.Content = strings.Join(contentLines, "\n")

	// meta data
	if err := metadata.Parse(item, lastModifiedDate, lines); err != nil {
		return fmt.Errorf("Unable to parse the meta data of item %q. Error: %s", item, err), nil
	}

	return
}
