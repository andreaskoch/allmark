// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"github.com/elWyatt/allmark/common/content"
	"github.com/elWyatt/allmark/common/route"
	"github.com/elWyatt/allmark/common/util/hashutil"
	"fmt"
)

type File struct {
	*content.ContentProvider

	parentRoute route.Route
	fileRoute   route.Route
}

func (file *File) String() string {
	return fmt.Sprintf("%s", file.fileRoute.Value())
}

func (file *File) Id() string {
	hash := hashutil.FromString(file.fileRoute.Value())
	return hash
}

func (file *File) Name() string {
	return fmt.Sprintf("%s", file.fileRoute.LastComponentName())
}

func (file *File) Parent() route.Route {
	return file.parentRoute
}

func (file *File) Route() route.Route {
	return file.fileRoute
}
