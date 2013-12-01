// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pattern

import (
	"fmt"
	"testing"
)

func Test_IsEmpty_NotEmpty(t *testing.T) {
	// arrange
	input := " asdas "

	// act
	result := IsEmpty(input)

	// assert
	if result != false {
		t.Errorf("The result should be false. The supplied string %q is NOT empty.", input)
	}
}

func Test_IsEmpty_EmptyString(t *testing.T) {
	// arrange
	input := ""

	// act
	result := IsEmpty(input)

	// assert
	if result != true {
		t.Errorf("The result should be true. The supplied string %q is empty.", input)
	}
}

func Test_IsEmpty_Süaces(t *testing.T) {
	// arrange
	input := "     "

	// act
	result := IsEmpty(input)

	// assert
	if result != true {
		t.Errorf("The result should be true. The supplied string %q is empty.", input)
	}
}

func Test_IsEmpty_SpacesAndTabs(t *testing.T) {
	// arrange
	input := "     \t   \t "

	// act
	result := IsEmpty(input)

	// assert
	if result != true {
		t.Errorf("The result should be true. The supplied string %q is empty.", input)
	}
}

func Test_IsEmpty_SpacesAndLineEndings(t *testing.T) {
	// arrange
	input := "     \t   \t \r\n \t   "

	// act
	result := IsEmpty(input)

	// assert
	if result != true {
		t.Errorf("The result should be true. The supplied string %q is empty.", input)
	}
}

func Test_IsHorizontalRule_ThreeDashes(t *testing.T) {
	// arrange
	input := "---"

	// act
	result := IsHorizontalRule(input)

	// assert
	if result != true {
		t.Errorf("The result should be true. The supplied string %q is a horizontal rule.", input)
	}
}

func Test_IsHorizontalRule_MultipleDashes(t *testing.T) {
	// arrange
	input := "--------"

	// act
	result := IsHorizontalRule(input)

	// assert
	if result != true {
		t.Errorf("The result should be true. The supplied string %q is a horizontal rule.", input)
	}
}

func Test_IsHorizontalRule_ThreeDashesWithWhitespaceBehind(t *testing.T) {
	// arrange
	input := "--- "

	// act
	result := IsHorizontalRule(input)

	// assert
	if result != true {
		t.Errorf("The result should be true. The supplied string %q is a horizontal rule.", input)
	}
}

func Test_IsHorizontalRule_ThreeDashesWithWhitespaceInFront(t *testing.T) {
	// arrange
	input := " ---"

	// act
	result := IsHorizontalRule(input)

	// assert
	if result != false {
		t.Errorf("The result should be false. The supplied string %q is NOT a horizontal rule.", input)
	}
}

func Test_IsHorizontalRule_TwoDashes(t *testing.T) {
	// arrange
	input := "--"

	// act
	result := IsHorizontalRule(input)

	// assert
	if result != false {
		t.Errorf("The result should be false. The supplied string %q is NOT a horizontal rule.", input)
	}
}

func Test_IsMetaDataDefinition_RandomText(t *testing.T) {
	// arrange
	input := "yada yada"

	// act
	result := IsMetaDataDefinition(input)

	// assert
	if result != false {
		t.Errorf("The result should be false. The supplied string %q is NOT a meta data defintion.", input)
	}
}

func Test_IsMetaDataDefinition_KeyOnly(t *testing.T) {
	// arrange
	input := "type:"

	// act
	result := IsMetaDataDefinition(input)

	// assert
	if result != true {
		t.Errorf("The result should be true. The supplied string %q is a meta data defintion.", input)
	}
}

func Test_IsMetaDataDefinition_KeyFollowedByWhitespace(t *testing.T) {
	// arrange
	input := "type:  "

	// act
	result := IsMetaDataDefinition(input)

	// assert
	if result != true {
		t.Errorf("The result should be true. The supplied string %q is a meta data defintion.", input)
	}
}

func Test_IsMetaDataDefinition_KeyValue(t *testing.T) {
	// arrange
	input := "type: document"

	// act
	result := IsMetaDataDefinition(input)

	// assert
	if result != true {
		t.Errorf("The result should be true. The supplied string %q is a meta data defintion.", input)
	}
}

func Test_GetMetaDataKey_KeyValue(t *testing.T) {
	// arrange
	key := "type"
	value := "document"
	input := fmt.Sprintf("%s: %s", key, value)

	// act
	result := GetMetaDataKey(input)

	// assert
	if result == "" {
		t.Errorf("The result should be %s but was %s.", key, result)
	}
}

func Test_GetMetaDataKey_KeyOnly(t *testing.T) {
	// arrange
	key := "type"
	input := fmt.Sprintf("%s:  ", key)

	// act
	result := GetMetaDataKey(input)

	// assert
	if result == "" {
		t.Errorf("The result should be %s but was %s.", key, result)
	}
}

func Test_GetMetaDataKey_NoLabel(t *testing.T) {
	// arrange
	input := "sadsa dasdasds"

	// act
	result := GetMetaDataKey(input)

	// assert
	if result != "" {
		t.Errorf("The result should be empty but was %s.", result)
	}
}

func Test_GetSingleLineMetaDataKeyAndValue_KeyValue(t *testing.T) {
	// arrange
	key := "type"
	value := "document"
	input := fmt.Sprintf("%s: %s", key, value)

	// act
	resultKey, resultValue := GetSingleLineMetaDataKeyAndValue(input)

	// assert
	if resultKey != key {
		t.Errorf("The result key should be %s but was %s.", key, resultKey)
	}

	if resultValue != value {
		t.Errorf("The result value should be %s but was %s.", key, resultValue)
	}
}

func Test_GetSingleLineMetaDataKeyAndValue_NoKeyValue(t *testing.T) {
	// arrange
	input := "Yada Yada"

	// act
	resultKey, resultValue := GetSingleLineMetaDataKeyAndValue(input)

	// assert
	if resultKey != "" {
		t.Errorf("The result key should be empty but was %s.", resultKey)
	}

	if resultValue != "" {
		t.Errorf("The result value should be empty but was %s.", resultValue)
	}
}
