// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

type Router interface {
	Route() *Route
}

type RootRouter struct {
	route *Route
}

func (rootRouter RootRouter) Route() *Route {
	if rootRouter.route == nil {
		rootRouter.route = New()
	}

	return rootRouter.route
}

func NewRootRouter() Router {
	return RootRouter{
		route: New(),
	}
}
