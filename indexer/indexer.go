// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package indexer

import (
	"andyk/docs/model"
	"fmt"
	"os"
	"path/filepath"
)

func walker(path string, _ os.FileInfo, _ error) error {
	fmt.Println(path)
	return nil
}

func Index(repositoryPath string) map[int]model.Document {

	docs := make(map[int]model.Document)

	// index the repository
	err := filepath.Walk(repositoryPath, walker)
	if err != nil {
		panic(err)
	}

	return docs
}
