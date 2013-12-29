// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package paths

type Pather interface {

	// Get the path for the supplied item
	Path(itemPath string) string
}
