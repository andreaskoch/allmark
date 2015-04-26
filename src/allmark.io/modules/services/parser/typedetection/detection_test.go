// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typedetection

import (
	"testing"

	"allmark.io/modules/model"
)

func Test_DetectType_Document(t *testing.T) {
	// arrange
	inputLines := []string{
		"",
		"---",
		"type: document",
	}
	expectedType := model.Document

	// act
	result := DetectType(inputLines)

	// assert
	if result != expectedType {
		t.Errorf("The result type should be %s but was %s", expectedType, result)
	}
}

func Test_DetectType_Presentation(t *testing.T) {
	// arrange
	inputLines := []string{
		"",
		"---",
		"type: presentation",
	}
	expectedType := model.Presentation

	// act
	result := DetectType(inputLines)

	// assert
	if result != expectedType {
		t.Errorf("The result type should be %s but was %s", expectedType, result)
	}
}

func Test_DetectType_Repository(t *testing.T) {
	// arrange
	inputLines := []string{
		"",
		"---",
		"type: repository",
	}
	expectedType := model.Repository

	// act
	result := DetectType(inputLines)

	// assert
	if result != expectedType {
		t.Errorf("The result type should be %s but was %s", expectedType, result)
	}
}

func Test_DetectType_UnknownTypeFallsBackToDocument(t *testing.T) {
	// arrange
	inputLines := []string{
		"",
		"---",
		"type: yada yada",
	}
	expectedType := model.Document

	// act
	result := DetectType(inputLines)

	// assert
	if result != expectedType {
		t.Errorf("The result type should be %s but was %s", expectedType, result)
	}
}
