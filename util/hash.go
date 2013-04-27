// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import (
	"crypto/sha1"
	"fmt"
)

func GetHash(text string) string {
	sha1 := sha1.New()
	sha1.Write([]byte(text))

	return fmt.Sprintf("%x", string(sha1.Sum(nil)[0:8]))
}
