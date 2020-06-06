// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metadata

import (
	"github.com/elWyatt/allmark/services/parser/pattern"
	"strings"
)

func GetMetaDataLines(lines []string) []string {

	lineNumber, err := GetMetaDataPosition(lines)
	if err != nil {
		return []string{} // there is no meta data in the supplied lines
	}

	// return the lines that contain the meta data
	if lineNumber+1 < len(lines) {
		return lines[(lineNumber + 1):]
	}

	// no meta data
	return []string{}
}

func getSingleLineMetaData(keyNames []string, lines []string) (keyFound bool, value string, remainingLines []string) {

	remainingLines = make([]string, 0)

	// search for the supplied key
	for _, line := range lines {
		if !keyFound {
			matchedKey, matchedValue := pattern.GetSingleLineMetaDataKeyAndValue(line)
			isKeyValuePair := !(matchedKey == "" && matchedValue == "")

			if isKeyValuePair && keyNamesMatch(matchedKey, keyNames) {
				value = strings.TrimSpace(matchedValue)
				keyFound = true
			}
		}

		remainingLines = append(remainingLines, line)
	}

	return keyFound, value, remainingLines
}

func keyNamesMatch(keyName string, keyNameAlternatives []string) bool {
	keyName = strings.ToLower(keyName)
	for _, keyNameAlternative := range keyNameAlternatives {
		if strings.ToLower(keyNameAlternative) == keyName {
			return true
		}
	}

	return false
}
