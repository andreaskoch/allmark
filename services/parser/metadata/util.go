// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metadata

import (
	"github.com/andreaskoch/allmark2/services/parser/pattern"
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
			key, value := pattern.GetSingleLineMetaDataKeyAndValue(line)
			isKeyValuePair := !(key == "" && value == "")

			if isKeyValuePair && keyNamesMatch(key, keyNames) {
				value = strings.TrimSpace(value)
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
