// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"github.com/andreaskoch/allmark/view"
)

type Message struct {
	Route     string     `json:"route"`
	Name      string     `json:"name"`
	ViewModel view.Model `json:"model"`
}

func UpdateMessage(viewModel view.Model) Message {
	return Message{
		Route:     viewModel.Route,
		Name:      "update",
		ViewModel: viewModel,
	}
}
