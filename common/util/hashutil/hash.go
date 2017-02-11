// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashutil

import (
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
)

func FromString(text string) string {
	return FromBytes([]byte(text))
}

func GetHash(reader io.Reader) (string, error) {

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return FromBytes(bytes), nil
}

func FromBytes(bytes []byte) string {
	crc := crc32.ChecksumIEEE(bytes)
	return fmt.Sprintf(`%d-%08X`, len(bytes), crc)
}
