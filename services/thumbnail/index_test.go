// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package thumbnail

import (
	"github.com/elWyatt/allmark/common/route"
	"testing"
)

func Test_GetThumbnailDimensionsFromRoute(t *testing.T) {
	// arrange
	requestURL := "/collections/Design/Splashscreens/collections/Design/Splashscreens/files/login-Space-Invaders.jpg-maxWidth:400-maxHeight:0"
	requestRoute := route.NewFromRequest(requestURL)

	// act
	resultRoute, _ := GetThumbnailDimensionsFromRoute(requestRoute)

	// assert
	expectedRoute := route.NewFromRequest("/collections/Design/Splashscreens/collections/Design/Splashscreens/files/login-Space-Invaders.jpg")
	if expectedRoute.Value() != resultRoute.Value() {
		t.Errorf("The base route should be %q but was %q.", expectedRoute, resultRoute)
	}
}
