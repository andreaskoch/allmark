// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package update

import (
	"allmark.io/modules/common/route"
	"allmark.io/modules/web/view/viewmodel"
)

type Message struct {
	Route       string           `json:"route"`
	Name        string           `json:"name"`
	UpdateModel viewmodel.Update `json:"model"`
}

func NewMessage(updateModel viewmodel.Update) Message {
	route := route.NewFromRequest(updateModel.Route)

	return Message{
		Route:       route.Value(),
		Name:        "update",
		UpdateModel: updateModel,
	}
}
