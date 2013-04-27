// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package view

import (
	"fmt"
)

type Model struct {
	Path        string
	Title       string
	Description string
	Content     string
	LanguageTag string
	Type        string
	Entries     []Model
}

func Error(msg string, path string) Model {
	return Model{
		Title:   fmt.Sprintf("Error: %s", msg),
		Path:    path,
		Content: msg,
		Type:    "error",
	}
}
