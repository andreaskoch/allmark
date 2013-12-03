// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pattern

import (
	"regexp"
	"strings"
)

var (
	// Lines which contain nothing but white space characters
	// or no characters at all.
	emptyLinePattern = regexp.MustCompile(`^\s*$`)

	// Lines which a start with a hash, followed by zero or more
	// white space characters, followed by text.
	titlePattern = regexp.MustCompile(`^#\s*([\pL\pN\p{Latin}]+.+)`)

	// Lines which start with text
	descriptionPattern = regexp.MustCompile(`^[\pL\pN\p{Latin}]+.+`)

	// Lines which nothing but dashes
	horizontalRulePattern = regexp.MustCompile(`^-{3,}\s*$`)

	// Lines with a "key: value" syntax
	singleLineMetaDataPattern = regexp.MustCompile(`^(\w+[\w\s]+\w+):\s*([\pL\pN\p{Latin}]+.+)$`)

	// Multi-line tags meta data
	multiLineTagsPattern = regexp.MustCompile(`(?is)tags:\n{0,2}(\n\s?-\s?[^\n]+)+\n*`)

	// Multi-line locations meta data
	multiLineLocationsPattern = regexp.MustCompile(`(?is)locations:\n{0,2}(\n\s?-\s?[^\n]+)+\n*`)

	// Lines with a meta data label in them syntax
	metaDataLabelPattern = regexp.MustCompile(`^(\w+[\w\s]+\w+):`)

	// Meta data list item pattern
	metaDataListItemPattern = regexp.MustCompile(`^\s?[*-]\s?(.+)$`)

	// Pattern which matches all HTML/XML tags
	HtmlTagPattern = regexp.MustCompile(`\<[^\>]*\>`)

	// Markdown headline pattern
	AnyLevelMarkdownHeadline = regexp.MustCompile(`^(#+?)([^#].+?[^#])(#*)$`)
)

func IsEmpty(line string) bool {
	return emptyLinePattern.MatchString(line)
}

func IsHorizontalRule(line string) bool {
	return horizontalRulePattern.MatchString(line)
}

func IsMetaDataDefinition(line string) bool {
	return metaDataLabelPattern.MatchString(line)
}

func GetMetaDataKey(line string) string {
	matches := metaDataLabelPattern.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

func GetSingleLineMetaDataKeyAndValue(line string) (key string, value string) {
	matches := singleLineMetaDataPattern.FindStringSubmatch(line)
	if len(matches) > 2 {
		return matches[1], matches[2]
	}

	if len(matches) > 1 {
		return matches[1], ""
	}

	return "", ""
}

func IsTitle(line string) (bool, string) {
	matches := titlePattern.FindStringSubmatch(line)
	if len(matches) > 1 {
		return true, matches[1] // title was found
	}

	return false, "" // no title was found
}

func IsDescription(line string) (bool, string) {
	if descriptionPattern.MatchString(line) {
		return true, line // description was found
	}

	return false, "" // no description was found
}

func IsListItem(line string) (bool, string) {
	matches := metaDataListItemPattern.FindStringSubmatch(line)
	if len(matches) > 1 {
		return true, strings.TrimSpace(matches[1]) // list item was found
	}

	return false, "" // no list item was found
}

func IsMultiLineTagDefinition(text string) (bool, []string) {
	return isMultiLineDefinition(multiLineTagsPattern, text)
}

func IsMultiLineLocationDefinition(text string) (bool, []string) {
	return isMultiLineDefinition(multiLineLocationsPattern, text)
}

func isMultiLineDefinition(pattern *regexp.Regexp, text string) (bool, []string) {
	multiLineTagLocation := pattern.FindStringSubmatchIndex(text)
	if multiLineTagLocation == nil {
		return false, []string{}
	}

	tagNames := make([]string, 0)

	multiLineTagBlock := strings.TrimSpace(text[multiLineTagLocation[0]:multiLineTagLocation[1]])
	tagLines := strings.Split(multiLineTagBlock, "\n")
	for _, line := range tagLines {
		if isListItem, tagName := IsListItem(line); isListItem {
			tagNames = append(tagNames, tagName)
		}
	}

	return true, tagNames
}
