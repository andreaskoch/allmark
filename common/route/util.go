// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

import (
	"strings"
)

func EncodeUrl(rawUrl string) string {
	encodedUrl := strings.Replace(rawUrl, "%", "%25", -1)
	encodedUrl = strings.Replace(encodedUrl, "+", "%2B", -1)
	encodedUrl = strings.Replace(encodedUrl, " ", "+", -1)
	encodedUrl = strings.Replace(encodedUrl, "#", "%23", -1)
	encodedUrl = strings.Replace(encodedUrl, "$", "%24", -1)
	encodedUrl = strings.Replace(encodedUrl, "&", "%26", -1)
	encodedUrl = strings.Replace(encodedUrl, ",", "%2C", -1)
	encodedUrl = strings.Replace(encodedUrl, ":", "%3A", -1)
	encodedUrl = strings.Replace(encodedUrl, ";", "%3B", -1)
	encodedUrl = strings.Replace(encodedUrl, "=", "%3D", -1)
	encodedUrl = strings.Replace(encodedUrl, "?", "%3F", -1)
	encodedUrl = strings.Replace(encodedUrl, "@", "%40", -1)
	return encodedUrl
}

func DecodeUrl(encodedUrl string) string {
	decodedUrl := strings.Replace(encodedUrl, "%25", "%", -1)
	decodedUrl = strings.Replace(decodedUrl, "%20", " ", -1)
	decodedUrl = strings.Replace(decodedUrl, "%23", "#", -1)
	decodedUrl = strings.Replace(decodedUrl, "%24", "$", -1)
	decodedUrl = strings.Replace(decodedUrl, "%26", "&", -1)
	decodedUrl = strings.Replace(decodedUrl, "+", " ", -1)
	decodedUrl = strings.Replace(decodedUrl, "%2B", "+", -1)
	decodedUrl = strings.Replace(decodedUrl, "%2C", ",", -1)
	decodedUrl = strings.Replace(decodedUrl, "%3A", ":", -1)
	decodedUrl = strings.Replace(decodedUrl, "%3B", ";", -1)
	decodedUrl = strings.Replace(decodedUrl, "%3D", "=", -1)
	decodedUrl = strings.Replace(decodedUrl, "%3F", "?", -1)
	decodedUrl = strings.Replace(decodedUrl, "%40", "@", -1)
	return decodedUrl
}
