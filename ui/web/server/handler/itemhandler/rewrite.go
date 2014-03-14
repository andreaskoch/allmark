// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package itemhandler

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/route"
	"regexp"
)

func NewRewrite(pattern, target string) RequestRewrite {
	rewritePattern := regexp.MustCompile(pattern)
	return RequestRewrite{rewritePattern, target}
}

type RequestRewrite struct {
	pattern *regexp.Regexp
	target  string
}

func (rewrite RequestRewrite) String() string {
	return fmt.Sprintf("Rewrite (%s → %s)", rewrite.pattern, rewrite.target)
}

func (rewrite RequestRewrite) Match(requestRoute route.Route) (bool, route.Route) {
	if !rewrite.pattern.MatchString(requestRoute.Value()) {
		return false, requestRoute
	}

	targetRoute, err := route.NewFromRequest(rewrite.pattern.ReplaceAllString(requestRoute.Value(), rewrite.target))
	if err != nil {
		panic(err)
	}

	return true, *targetRoute
}
