// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashutil

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"
)

func GetHash(reader io.Reader) (string, error) {

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	sha1Hash := sha1.New()
	sha1Hash.Write(bytes)
	hashBytes := sha1Hash.Sum(nil)

	return string(hex.EncodeToString(hashBytes)), nil
}
