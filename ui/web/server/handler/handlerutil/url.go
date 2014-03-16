// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlerutil

import (
	"net/url"
	"strconv"
)

func GetPageParameterFromUrl(url url.URL) (page int, parameterIsAvailable bool) {
	pageParam := url.Query().Get("page")
	if pageParam == "" {
		return 0, false
	}

	page64, err := strconv.ParseInt(pageParam, 10, 64)
	if err != nil {
		return 0, true
	}

	if page64 < 1 {
		return 0, true
	}

	return int(page64), true
}
