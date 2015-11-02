// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

type BreadcrumbNavigation struct {
	Entries []Breadcrumb `json:"entries"`
}

func (navigation BreadcrumbNavigation) IsAvailable() bool {
	return len(navigation.Entries) > 0
}

type Breadcrumb struct {
	Level  int    `json:"level"`
	Title  string `json:"title"`
	Path   string `json:"path"`
	IsLast bool
}
