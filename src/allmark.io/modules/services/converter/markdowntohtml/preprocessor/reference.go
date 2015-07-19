// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package preprocessor

import (
	"allmark.io/modules/common/paths"
	"allmark.io/modules/common/pattern"
	"allmark.io/modules/model"
	"fmt"
	"regexp"
	"strings"
)

var (
	// [reference:*alias-of-referenced-item*]
	referencePattern = regexp.MustCompile(`\[reference:([^\]]+)\]`)
)

func newReferenceExtension(pathProvider paths.Pather, aliasResolver func(alias string) *model.Item) *referenceExtension {
	return &referenceExtension{
		pathProvider:  pathProvider,
		aliasResolver: aliasResolver,
	}
}

type referenceExtension struct {
	pathProvider  paths.Pather
	aliasResolver func(alias string) *model.Item
}

func (converter *referenceExtension) Convert(markdown string) (convertedContent string, converterError error) {

	convertedContent = markdown

	for {

		// search for references
		found, matches := pattern.IsMatch(convertedContent, referencePattern)
		if !found || (found && len(matches) != 2) {
			break // abort. no (more) references found
		}

		// extract the parameters from the pattern matches
		originalText := strings.TrimSpace(matches[0])
		alias := strings.TrimSpace(matches[1])

		// lookup the item
		item := converter.aliasResolver(alias)
		if item == nil {
			// an item with the alias was not found
			convertedContent = strings.Replace(convertedContent, originalText, fmt.Sprintf("<!-- Alias %q not found -->", alias), 1)
			continue
		}

		// normalize the path with the current path provider
		path := converter.pathProvider.Path(item.Route().Value())

		// assemble the link
		linkCode := fmt.Sprintf("[%s](%s)", item.Title, path)

		// replace markdown with link list
		convertedContent = strings.Replace(convertedContent, originalText, linkCode, 1)

	}

	return convertedContent, nil
}
