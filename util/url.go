// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import (
	"net/url"
)

func EncodeUrl(rawurl string) string {
	url, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}

	return url.String()
}
