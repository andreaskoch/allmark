// Copyright 2014 Andreas Koch. All rights reserved.
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

func Test_GetSingleLineMetaDataKeyAndValue_KeyMultipleValues(t *testing.T) {
	// arrange
	key := "tags"
	value := "tag1, tag2, tag3"
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

func Test_IsMultiLineTagDefinition_True(t *testing.T) {
	// arrange
	expected := true
	input := `tags:
- tag1
- tag2
- tag3`

	// act
	result, _ := IsMultiLineTagDefinition(input)

	// assert
	if result != expected {
		t.Errorf("The text %q is a multi-line tag definition, but the result was %v (expected: %v).", input, result, expected)
	}
}

func Test_IsMultiLineTagDefinition_AllTagsAreReturned(t *testing.T) {
	// arrange
	expected := 3
	input := `tags:
- tag1
- tag2
- tag3`

	// act
	_, matches := IsMultiLineTagDefinition(input)

	// assert
	if len(matches) != expected {
		t.Errorf("The text %q should return %d matches but returned %d. Result: %#v", input, expected, len(matches), matches)
	}
}

func Test_IsDescription_TextContainsNoMarkdown_ResultIsTrue(t *testing.T) {
	// arrange
	input := "Yada Yada Some Text that does not contain markdown."
	expected := true

	// act
	result, _ := IsDescription(input)

	// assert
	if result != expected {
		t.Errorf("The text %q is not a description because it contains markdown but the IsDescription function returned %v (expected: %v).", input, result, expected)
	}
}

func Test_IsDescription_TextContainsMarkdown_ResultIsFalse(t *testing.T) {
	// arrange
	input := "Yada Yada Some Text that looks like a description but suddenly contains an ![image](files/image.png). And then some more text."
	expected := false

	// act
	result, _ := IsDescription(input)

	// assert
	if result != expected {
		t.Errorf("The text %q is not a description because it contains markdown but the IsDescription function returned %v (expected: %v).", input, result, expected)
	}
}
