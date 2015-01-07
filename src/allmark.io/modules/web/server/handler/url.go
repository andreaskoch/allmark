// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"net/url"
	"strconv"
)

func getPageParameterFromUrl(url url.URL) (page int, parameterIsAvailable bool) {
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

func getQueryParameterFromUrl(url url.URL) (query string, parameterIsAvailable bool) {
	queryParam := url.Query().Get("q")
	if queryParam == "" {
		return "", false
	}

	return queryParam, true
}
