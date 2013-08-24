// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pattern

import (
	"regexp"
)

var (
	// Lines which contain nothing but white space characters
	// or no characters at all.
	EmptyLinePattern = regexp.MustCompile(`^\s*$`)

	// Lines which a start with a hash, followed by zero or more
	// white space characters, followed by text.
	TitlePattern = regexp.MustCompile(`^#\s*([\pL\pN\p{Latin}]+.+)`)

	// Lines which start with text
	DescriptionPattern = regexp.MustCompile(`^[\pL\pN\p{Latin}]+.+`)

	// Lines which nothing but dashes
	HorizontalRulePattern = regexp.MustCompile(`^-{2,}`)

	// Lines with a "key: value" syntax
	MetaDataPattern = regexp.MustCompile(`^(\w+):\s*([\pL\pN\p{Latin}]+.+)$`)

	// Pattern which matches all HTML/XML tags
	HtmlTagPattern = regexp.MustCompile(`\<[^\>]*\>`)
)
