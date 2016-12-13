// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package postprocessor

import (
	"github.com/kyokomi/emoji"
	"regexp"
	"strings"
)

var (
	emojiPattern = regexp.MustCompile(`:[\w\d_]+:`)
)

// addEmojis searches the supplied HTML code for supported emojis
// and replaces them with the respective emoji icon (see: http://www.emoji-cheat-sheet.com/).
// Example: :dancers: becomes 👯
func addEmojis(html string) string {

	allMatches := emojiPattern.FindAllStringSubmatch(html, -1)
	for _, matches := range allMatches {
		if len(matches) == 0 {
			continue
		}

		originalText := strings.TrimSpace(matches[0])
		emojifiedText := emoji.Sprint(originalText)

		// replace all occurances
		html = strings.Replace(html, originalText, emojifiedText, -1)

	}

	return emoji.Sprint(html)
}
