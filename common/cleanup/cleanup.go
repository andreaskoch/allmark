// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cleanup

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/shutdown"
	"os"
	"strings"
	"time"
)

var removeNow = make(map[string]bool)
var removeLater = make(map[string]bool)
var isRunning = true

func init() {

	shutdown.Register(func() error {

		// stop all running processes
		isRunning = false

		// trigger the cleanup
		return cleanup()

	})

	instantRemovalProcess()
}

func Now(filename string) {
	removeNow[filename] = true
}

func OnShutdown(filename string) {
	removeLater[filename] = true
}

// Try to remove all files which have been queed for deletion.
func cleanup() error {

	// todo: debug message
	fmt.Printf("Things to remove %v\n", removeLater)

	errors := make([]string, 0)
	for entry, _ := range removeLater {

		// todo: debug message
		fmt.Printf("Removing %q\n", entry)
		if err := os.RemoveAll(entry); err != nil && os.IsNotExist(err) == false {

			// todo: find a better way to aggregate errors
			errors = append(errors, fmt.Sprintf("Failed to remove %q \n", entry))
		}

	}

	if len(errors) == 0 {
		return nil
	}

	return fmt.Errorf("Cleanup errors: %s", strings.Join(errors, "\n"))
}

func instantRemovalProcess() {
	go func() {
		for isRunning {

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

			if isRunning {
				time.Sleep(time.Second * 2) // wait a few seconds before the next run
			}
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
