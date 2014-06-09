// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updatehandler

import (
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
)

type Message struct {
	Route     string          `json:"route"`
	Name      string          `json:"name"`
	ViewModel viewmodel.Model `json:"model"`
}

func UpdateMessage(viewModel viewmodel.Model) Message {
	return Message{
		Route:     viewModel.Route,
		Name:      "update",
		ViewModel: viewModel,
	}
}
