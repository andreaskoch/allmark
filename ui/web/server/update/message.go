// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package update

import (
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
)

type Message struct {
	Route     string          `json:"route"`
	Name      string          `json:"name"`
	ViewModel viewmodel.Model `json:"model"`
}

func NewMessage(viewModel viewmodel.Model) Message {
	route, err := route.NewFromRequest(viewModel.Route)
	if err != nil {
		panic(err)
	}

	return Message{
		Route:     route.Value(),
		Name:      "update",
		ViewModel: viewModel,
	}
}
