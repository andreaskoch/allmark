// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package files

import (
	"strings"
)

func getLastPathComponent(path string) string {
	if !strings.Contains(path, "/") {
		return path
	}

	components := strings.Split(path, "/")
	return components[len(components)-1]
}
