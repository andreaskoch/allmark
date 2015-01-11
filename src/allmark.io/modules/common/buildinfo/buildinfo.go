// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package buildinfo provides information about the current build.
package buildinfo

// GitInfo is either the empty string (the default)
// or is set to the git hash of the most recent commit
// using the -X linker flag (Example: "2015-01-11-284c030+")
var GitInfo string

// Version returns the git version of this binary (e.g. "2015-01-11-284c030+").
// If the linker flags were not provided, the return value is "unknown".
func Version() string {
	if GitInfo != "" {
		return GitInfo
	}
	return "unknown"
}
