// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/dataaccess/filesystem"
)

func main() {

	repository, err := filesystem.NewRepository(fsutil.GetWorkingDirectory())
	if err != nil {
		panic(err)
	}

	itemEvents, done := repository.GetItems()

	allItemsRetrieved := false
	for !allItemsRetrieved {
		select {
		case allItemsRetrieved = <-done:
		case itemEvent := <-itemEvents:
			if itemEvent != nil {
				fmt.Println(itemEvent.Item)
			}
		}
	}
}
