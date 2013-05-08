// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package view

import (
	"fmt"
)

type Model struct {
	Route       string
	Title       string
	Description string
	Content     string
	LanguageTag string
	Type        string
	Date        string
}

func Error(msg string, path string) Model {
	return Model{
		Title:   fmt.Sprintf("Error: %s", msg),
		Route:   path,
		Content: msg,
		Type:    "error",
	}
}
