// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	util "github.com/andreaskoch/allmark2/common/util/filesystem"
	"github.com/andreaskoch/allmark2/dataaccess"
	"github.com/andreaskoch/allmark2/dataaccess/filesystem"
)

func main() {

	filesystemAccessor, err := filesystem.New(util.GetWorkingDirectory())
	if err != nil {
		panic(err)
	}

	rootItem, err := filesystemAccessor.GetRootItem()
	if err != nil {
		panic(err)
	}

	rootItem.Walk(func(item *dataaccess.Item) {
		fmt.Println(item)
	})
}
