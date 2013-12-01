// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package document

import (
	"fmt"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/parser/pattern"
	"strings"
)

func Parse(item *model.Item, lines []string) error {

	// title
	titleLineNumber := len(lines)
	for lineNumber, line := range lines {

		// ignore empty lines
		if pattern.IsEmpty(line) {
			continue
		}

		// search for the title
		isTitle, title := pattern.IsTitle(line)
		if !isTitle {
			return fmt.Errorf("The line %q does not contain a title.", line)
		}

		// capture the line number of the title
		titleLineNumber = lineNumber

		// save the title to the item
		item.Title = strings.TrimSpace(title)
		break
	}

	// abort if there are no more lines
	if len(lines) < titleLineNumber+1 {
		return nil // there are no more lines, but having no description is ok
	}

	// description
	descriptionLineNumber := len(lines)
	for lineNumber, line := range lines[(titleLineNumber + 1):] {

		// ignore empty lines
		if pattern.IsEmpty(line) {
			continue
		}

		// search for the description
		isDescription, description := pattern.IsDescription(line)
		if !isDescription {
			return fmt.Errorf("The line %q does not contain a description.", line)
		}

		// capture the line number of the description
		descriptionLineNumber = lineNumber

		// save the description to the item
		item.Description = strings.TrimSpace(description)
		break
	}

	// abort if there are no more lines
	if len(lines) < descriptionLineNumber+1 {
		return nil // there are no more lines, but having no content is ok
	}

	// content

	return nil
}
