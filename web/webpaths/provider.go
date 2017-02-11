// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webpaths

import (
	"github.com/andreaskoch/allmark/common/paths"
	"github.com/andreaskoch/allmark/common/route"
)

type WebPathProvider struct {
	patherFactory paths.PatherFactory
	itemPather    paths.Pather
	tagPather     paths.Pather
}

func NewWebPathProvider(patherFactory paths.PatherFactory, basePath, tagPathPrefix string) WebPathProvider {
	return WebPathProvider{
		patherFactory: patherFactory,
		itemPather:    patherFactory.Absolute(basePath),
		tagPather:     patherFactory.Absolute(tagPathPrefix),
	}
}

func (provider *WebPathProvider) AbsolutePather(prefix string) paths.Pather {
	return provider.patherFactory.Absolute(prefix)
}

func (provider *WebPathProvider) ItemPather() paths.Pather {
	return provider.itemPather
}

func (provider *WebPathProvider) TagPather() paths.Pather {
	return provider.tagPather
}

func (provider *WebPathProvider) RelativePather(baseRoute route.Route) paths.Pather {
	return provider.patherFactory.Relative(baseRoute)
}
