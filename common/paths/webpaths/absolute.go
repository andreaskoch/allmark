// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webpaths

// Create a new absolute web path provider
func newAbsoluteWebPathProvider(base string) *AbsoluteWebPathProvider {
	return &AbsoluteWebPathProvider{
		base: base,
	}
}

type AbsoluteWebPathProvider struct {
	base string
}

// Get the absolute path for the supplied item
func (webPathProvider *AbsoluteWebPathProvider) Path(itemPath string) string {
	return webPathProvider.base + "/" + itemPath
}
