// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package defaulttheme

var (
	templates = make(map[string]string)
)

// RawTemplates returns a map of all raw templates in this theme by their name.
func RawTemplates() map[string]string {
	return templates
}
