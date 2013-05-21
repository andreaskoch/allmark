// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package view

import (
	"fmt"
)

type Model struct {
	Level         int      `json:"level"`
	AbsoluteRoute string   `json:"absoluteRoute"`
	RelativeRoute string   `json:"relativeRoute"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Content       string   `json:"content"`
	LanguageTag   string   `json:"languageTag"`
	Type          string   `json:"type"`
	Date          string   `json:"date"`
	SubEntries    []*Model `json:"subEntries"`
}

func Error(msg, relativPath, absolutePath string) *Model {
	return &Model{
		Level:         0,
		Title:         fmt.Sprintf("Error: %s", msg),
		RelativeRoute: relativPath,
		AbsoluteRoute: absolutePath,
		Content:       msg,
		Type:          "error",
	}
}
