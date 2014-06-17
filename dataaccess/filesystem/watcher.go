// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"fmt"
	"github.com/andreaskoch/go-fswatch"
)

var watchers map[string]bool

func init() {
	watchers = make(map[string]bool)
}

func watchFolder(folder string, callback func(change *fswatch.FolderChange)) {

	if exists, _ := watchers[folder]; exists {
		fmt.Println("Watcher already exists")
		return
	}

	// look for changes in the item directory
	go func() {
		var skipFunc = func(path string) bool {
			isReserved := isReservedDirectory(path)
			return isReserved
		}

		recurse := false
		checkIntervalInSeconds := 3
		folderWatcher := fswatch.NewFolderWatcher(folder, recurse, skipFunc, checkIntervalInSeconds).Start()

		for folderWatcher.IsRunning() {

			select {
			case change := <-folderWatcher.Change:

				callback(change)

			}

		}

		delete(watchers, folder)
		fmt.Printf("Exiting directory listener for folder %q.\n", folder)
	}()

	watchers[folder] = true
}
