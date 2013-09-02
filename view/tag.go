// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package view

type TagMap struct {
	Tags []*Tag
}

type Tag struct {
	Name          string   `json:"name"`
	AbsoluteRoute string   `json:"absoluteRoute"`
	Description   string   `json:"description"`
	Childs        []*Model `json:"childs"`
}
