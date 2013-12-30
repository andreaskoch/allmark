// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

type BreadcrumbNavigation struct {
	Entries []*Breadcrumb
}

func (navigation *BreadcrumbNavigation) IsAvailable() bool {
	return len(navigation.Entries) > 0
}

type Breadcrumb struct {
	Level int
	Title string
	Path  string
}
