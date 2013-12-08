// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typedetection

import (
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/parsing/metadata"
	"github.com/andreaskoch/allmark2/services/parsing/pattern"
	"strings"
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
	case "message":
		return model.TypeMessage
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
