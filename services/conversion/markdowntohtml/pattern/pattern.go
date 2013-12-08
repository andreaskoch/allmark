// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pattern

import (
	"regexp"
)

// IsMatch returns a flag indicating whether the supplied
// text and pattern do match and if yet, the matched text.
func IsMatch(text string, pattern *regexp.Regexp) (isMatch bool, matches []string) {
	matches = pattern.FindStringSubmatch(text)
	return matches != nil, matches
}
