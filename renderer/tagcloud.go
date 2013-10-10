// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renderer

import (
	"github.com/andreaskoch/allmark/mapper"
	"github.com/andreaskoch/allmark/repository"
)

func attachTagCloud(item *repository.Item) {
	cloud := mapper.MapTagCloud(tags)
	item.Model.TagCloud = &cloud
}
