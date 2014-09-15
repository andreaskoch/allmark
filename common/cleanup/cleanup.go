// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cleanup

import (
	"fmt"
	"os"
	"time"
)

var instantCleanup chan string
var removeNow map[string]bool

var shutdownCleanup chan string
var removeLater map[string]bool

func init() {
	instantCleanup = make(chan string, 1)
	removeNow = make(map[string]bool)

	shutdownCleanup = make(chan string, 1)
	removeLater = make(map[string]bool)

	go func() {
		for {
			select {

			case filename := <-instantCleanup:
				{
					removeNow[filename] = true
				}

			case filename := <-shutdownCleanup:
				{
					removeLater[filename] = true
				}

			}
		}
	}()

	instantRemovalProcess()
}

func Now(filename string) {
	go func() {
		instantCleanup <- filename
	}()
}

func OnShutdown(filename string) {
	go func() {
		shutdownCleanup <- filename
	}()
}

func Cleanup() {

	// try to remove all files from the list
	for filename, _ := range removeLater {
		if err := os.Remove(filename); err != nil && os.IsNotExist(err) == false {
			fmt.Printf("Failed to remove %q \n", filename)
		}
	}

}

func instantRemovalProcess() {
	go func() {
		for {

			// try to remove all files from the list
			for _, filename := range getFilesFromMap(removeNow) {

				if err := os.Remove(filename); err == nil || os.IsNotExist(err) {

					// file was successfully deleted
					delete(removeNow, filename)

					continue
				}

				// remember the file for the final cleanup
				OnShutdown(filename)
			}

			// wait a few seconds before the next run
			time.Sleep(time.Second * 2)
		}
	}()
}

func getFilesFromMap(fileMap map[string]bool) []string {
	files := make([]string, 0)
	for filename, _ := range fileMap {
		files = append(files, filename)
	}
	return files
}
