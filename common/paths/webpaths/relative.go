// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webpaths

// Create a new relative web path provider
func newRelativeWebPathProvider() *RelativeWebPathProvider {
	return &RelativeWebPathProvider{}
}

type RelativeWebPathProvider struct {
}

// Get the path relative for the supplied item
func (webPathProvider *RelativeWebPathProvider) Path(itemPath string) string {
	return itemPath
}
