// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"github.com/andreaskoch/allmark/util"
	"regexp"
	"strings"
	"time"
)

var (
	iso6391TwoLetterLanguageCodePattern = regexp.MustCompile(`^[a-z]$`)
	ietfLanguageTagPattern              = regexp.MustCompile(`^(\w\w)-\w{2,3}$`)
)

// Get ISO 639-1 language code from a given language string (e.g. "en-US" => "en", "de-DE" => "de")
func getTwoLetterLanguageCode(languageString string) string {

	fallbackLangueCode := "en"
	if languageString == "" {
		// default value
		return fallbackLangueCode
	}

	// Check if the language string already matches
	// the ISO 639-1 language code pattern (e.g. "en", "de").
	if len(languageString) == 2 && iso6391TwoLetterLanguageCodePattern.MatchString(languageString) {
		return strings.ToLower(languageString)
	}

	// Check if the language string matches the
	// IETF language tag pattern (e.g. "en-US", "de-DE").
	matchesIETFPattern, matches := util.IsMatch(languageString, ietfLanguageTagPattern)
	if matchesIETFPattern {
		return matches[1]
	}

	// use fallback
	return fallbackLangueCode
}

func formatDate(date time.Time) string {
	return date.Format("2006-01-02")
}
