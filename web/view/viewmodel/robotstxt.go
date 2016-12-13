// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

// RobotsTxtDisallow is the definition of disallow-section of a robots.txt
type RobotsTxtDisallow struct {
	UserAgent string
	Paths     []string
}

// RobotsTxt represents the content of a robots.txt file
type RobotsTxt struct {
	Disallows  []RobotsTxtDisallow
	SitemapURL string
}
