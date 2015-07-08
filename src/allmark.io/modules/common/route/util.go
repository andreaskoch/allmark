// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

import (
	"strings"
)

func ToKey(route Route) string {
	return route.Value()
}

func EncodeURL(rawURL string) string {
	encodedURL := strings.Replace(rawURL, "%", "%25", -1)
	encodedURL = strings.Replace(encodedURL, "+", "%2B", -1)
	encodedURL = strings.Replace(encodedURL, " ", "+", -1)
	encodedURL = strings.Replace(encodedURL, "#", "%23", -1)
	encodedURL = strings.Replace(encodedURL, "$", "%24", -1)
	encodedURL = strings.Replace(encodedURL, "&", "%26", -1)
	encodedURL = strings.Replace(encodedURL, ",", "%2C", -1)
	encodedURL = strings.Replace(encodedURL, ":", "%3A", -1)
	encodedURL = strings.Replace(encodedURL, ";", "%3B", -1)
	encodedURL = strings.Replace(encodedURL, "=", "%3D", -1)
	encodedURL = strings.Replace(encodedURL, "?", "%3F", -1)
	encodedURL = strings.Replace(encodedURL, "@", "%40", -1)
	return encodedURL
}

func DecodeURL(encodedURL string) string {
	decodedURL := strings.Replace(encodedURL, "%25", "%", -1)
	decodedURL = strings.Replace(decodedURL, "%20", " ", -1)
	decodedURL = strings.Replace(decodedURL, "%23", "#", -1)
	decodedURL = strings.Replace(decodedURL, "%24", "$", -1)
	decodedURL = strings.Replace(decodedURL, "%26", "&", -1)
	decodedURL = strings.Replace(decodedURL, "+", " ", -1)
	decodedURL = strings.Replace(decodedURL, "%2B", "+", -1)
	decodedURL = strings.Replace(decodedURL, "%2C", ",", -1)
	decodedURL = strings.Replace(decodedURL, "%3A", ":", -1)
	decodedURL = strings.Replace(decodedURL, "%3B", ";", -1)
	decodedURL = strings.Replace(decodedURL, "%3D", "=", -1)
	decodedURL = strings.Replace(decodedURL, "%3F", "?", -1)
	decodedURL = strings.Replace(decodedURL, "%40", "@", -1)
	return decodedURL
}
