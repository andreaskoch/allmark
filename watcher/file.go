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
	Modified chan bool
	Moved    chan bool

	debug   bool
	file    string
	running bool
}

func NewFileWatcher(filePath string) *FileWatcher {
	return &FileWatcher{
		Modified: make(chan bool),
		Moved:    make(chan bool),
		debug:    false,
		file:     filePath,
	}
}

func (fileWatcher *FileWatcher) String() string {
	return fmt.Sprintf("Filewatcher %q", fileWatcher.file)
}

func (fileWatcher *FileWatcher) Start() *FileWatcher {
	fileWatcher.running = true
	sleepTime := time.Second * 2

	go func() {

		for fileWatcher.running {

			if fileInfo, err := os.Stat(fileWatcher.file); err == nil {

				// check if file has been modified
				sleepTime := time.Now().Add(sleepTime * -1)
				modTime := fileInfo.ModTime()
				if sleepTime.Before(modTime) {

					// send out the notification
					fileWatcher.log("Item was modified")
					go func() {
						fileWatcher.Modified <- true
					}()
				}

			} else if os.IsNotExist(err) {

				// send out the notification
				fileWatcher.log("Item was removed")
				go func() {
					fileWatcher.Moved <- true
				}()

				// stop this file watcher
				fileWatcher.Stop()
			}

			time.Sleep(sleepTime)

		}

		fileWatcher.log("Stopped")
	}()

	return fileWatcher
}

func (fileWatcher *FileWatcher) Stop() *FileWatcher {
	fileWatcher.log("Stopping")
	fileWatcher.running = false
	return fileWatcher
}

func (fileWatcher *FileWatcher) IsRunning() bool {
	return fileWatcher.running
}

func (fileWatcher *FileWatcher) log(message string) *FileWatcher {
	if fileWatcher.debug {
		fmt.Printf("%s - %s\n", fileWatcher, message)
	}

	return fileWatcher
}
