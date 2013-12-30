// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

type Sitemap struct {
	AbsoluteRoute string     `json:"absoluteRoute"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Childs        []*Sitemap `json:"childs"`
}
