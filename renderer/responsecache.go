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

func cacheReponse(targetFile string, pathProvider *path.Provider, responseWriter ResponseWriter, host string, writer io.Writer) {

	// assemble the file path on disk
	renderTargetPath := pathProvider.GetRenderTargetPathForGiven(targetFile)

	// read the file from disk if it exists
	if util.FileExists(renderTargetPath) {
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

	defer file.Close()

	fileWriter := bufio.NewWriter(file)
	responseWriter(fileWriter, host)

	fileWriter.Flush()

	// try to read from the cache again
	cacheReponse(targetFile, pathProvider, responseWriter, host, writer)
}
