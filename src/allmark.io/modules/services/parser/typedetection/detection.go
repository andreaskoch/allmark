// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typedetection

import (
	"strings"

	"allmark.io/modules/model"
	"allmark.io/modules/services/parser/metadata"
	"allmark.io/modules/services/parser/pattern"
)

func DetectType(lines []string) model.ItemType {

	// get the meta data definitions
	lines = metadata.GetMetaDataLines(lines)
	if len(lines) == 0 {
		return model.TypeDocument
	}

	// find the type name
	typeName := ""
	for _, line := range lines {
		if !pattern.IsMetaDataDefinition(line) {
			continue
		}

		// search for a type definition
		if key, value := pattern.GetSingleLineMetaDataKeyAndValue(line); strings.ToLower(key) == "type" && strings.TrimSpace(value) != "" {
			typeName = strings.TrimSpace(strings.ToLower(value))
			break // found a type definition
		}
	}

	// detect the type
	switch typeName {
	case "document":
		return model.TypeDocument
	case "location":
		return model.TypeLocation
	case "presentation":
		return model.TypePresentation
	case "repository":
		return model.TypeRepository
	default:
		return model.TypeDocument // fallback
	}

	panic("Unreachable")
}
