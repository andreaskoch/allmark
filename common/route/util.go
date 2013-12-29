// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

import (
	"strings"
)

func EncodeUrl(rawUrl string) string {
	rawUrl = strings.Replace(rawUrl, "%", "%25", -1)
	rawUrl = strings.Replace(rawUrl, " ", "%20", -1)
	rawUrl = strings.Replace(rawUrl, "#", "%23", -1)
	rawUrl = strings.Replace(rawUrl, "$", "%24", -1)
	rawUrl = strings.Replace(rawUrl, "&", "%26", -1)
	rawUrl = strings.Replace(rawUrl, "+", "%2B", -1)
	rawUrl = strings.Replace(rawUrl, ",", "%2C", -1)
	rawUrl = strings.Replace(rawUrl, ":", "%3A", -1)
	rawUrl = strings.Replace(rawUrl, ";", "%3B", -1)
	rawUrl = strings.Replace(rawUrl, "=", "%3D", -1)
	rawUrl = strings.Replace(rawUrl, "?", "%3F", -1)
	rawUrl = strings.Replace(rawUrl, "@", "%40", -1)
	return rawUrl
}

func DecodeUrl(encodedUrl string) string {
	encodedUrl = strings.Replace(encodedUrl, "%25", "%", -1)
	encodedUrl = strings.Replace(encodedUrl, "%20", " ", -1)
	encodedUrl = strings.Replace(encodedUrl, "%23", "#", -1)
	encodedUrl = strings.Replace(encodedUrl, "%24", "$", -1)
	encodedUrl = strings.Replace(encodedUrl, "%26", "&", -1)
	encodedUrl = strings.Replace(encodedUrl, "%2B", "+", -1)
	encodedUrl = strings.Replace(encodedUrl, "%2C", ",", -1)
	encodedUrl = strings.Replace(encodedUrl, "%3A", ":", -1)
	encodedUrl = strings.Replace(encodedUrl, "%3B", ";", -1)
	encodedUrl = strings.Replace(encodedUrl, "%3D", "=", -1)
	encodedUrl = strings.Replace(encodedUrl, "%3F", "?", -1)
	encodedUrl = strings.Replace(encodedUrl, "%40", "@", -1)
	return encodedUrl
}
