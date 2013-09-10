// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renderer

import (
	"bufio"
	"fmt"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/util"
	"io"
	"os"
)

var cachedResponseFiles map[string]bool

func init() {
	cachedResponseFiles = make(map[string]bool)
}

func cacheReponse(targetFile string, pathProvider *path.Provider, responseWriter ResponseWriter, host string, writer io.Writer) {

	// assemble the file path on disk
	renderTargetPath := pathProvider.GetRenderTargetPathForGiven(fmt.Sprintf("host-%s-%s", util.GetHash(host), targetFile))

	// read the file from disk if it exists
	isCached, renderTargetPathHasBeenCachedOnce := cachedResponseFiles[renderTargetPath]
	readResposeFromDisk := isCached && renderTargetPathHasBeenCachedOnce

	if util.FileExists(renderTargetPath) && readResposeFromDisk {
		file, err := os.Open(renderTargetPath)
		if err != nil {
			fmt.Fprintln(writer, fmt.Sprintf(`Unable to read %q.`, renderTargetPath))
			return
		}

		defer file.Close()

		data := make([]byte, 100)
		reader := bufio.NewReader(file)
		for {
			n, err := reader.Read(data)
			if n == 0 && err == io.EOF {
				return
			}

			if err != nil && err != io.EOF {
				fmt.Fprintln(writer, fmt.Sprintf(`Unable to read %q.`, renderTargetPath))
				return
			}

			writer.Write(data[:n])
		}

		return
	}

	// write to the file
	file, err := os.OpenFile(renderTargetPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Fprintln(writer, fmt.Sprintf(`Unable to open %q for writing. Err: %s`, renderTargetPath, err))
		return
	}

	defer func() {
		file.Close()

		// mark the target file as cached
		cachedResponseFiles[renderTargetPath] = true

		// try to read from the cache again
		cacheReponse(targetFile, pathProvider, responseWriter, host, writer)
	}()

	fileWriter := bufio.NewWriter(file)
	responseWriter(fileWriter, host)
	fileWriter.Flush()
}

func clearCachedResponses() {
	for filepath, isCached := range cachedResponseFiles {
		if !isCached {
			continue
		}

		if !util.FileExists(filepath) {
			cachedResponseFiles[filepath] = false
			continue
		}

		fmt.Printf("Removing cached response %q from disk.\n", filepath)
		if err := os.Remove(filepath); err != nil {
			fmt.Printf("Unable to remove file %q. Error: %s\n", filepath, err)
			continue
		}

		cachedResponseFiles[filepath] = false
	}
}
