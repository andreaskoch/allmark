// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watcher

import (
	"fmt"
	"os"
	"time"
)

type FileWatcher struct {
	file    string
	stopped chan bool
}

func NewFileWatcher(filePath string) *FileWatcher {
	return &FileWatcher{
		file: filePath,
	}
}

func (fileWatcher *FileWatcher) Start() {
	go func() {
		running := true
		sleepTime := time.Second * 2

		for running {
			select {
			default:

				if fileInfo, err := os.Stat(fileWatcher.file); err == nil {

					// check if file has been modified
					sleepTime := time.Now().Add(sleepTime * -1)
					modTime := fileInfo.ModTime()
					if sleepTime.Before(modTime) {
						fmt.Println("Item was modified")
					}

				} else if os.IsNotExist(err) {

					// file has been moved. check if it has been deleted
					// or if it has been renamed
					fmt.Println("Item was removed")
					running = false
				}

				time.Sleep(sleepTime)
			}
		}

		fmt.Println("Stopped")
	}()
}
