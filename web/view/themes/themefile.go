// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package themes

import (
	"encoding/base64"
	"fmt"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"io/ioutil"
	"path/filepath"
)

func newFileFromText(filename, text string) *ThemeFile {
	return &ThemeFile{
		path: filename,
		data: []byte(text),
	}
}

func newFileFromBase64(filename, base64Text string) *ThemeFile {

	data, err := base64.StdEncoding.DecodeString(base64Text)
	if err != nil {
		panic(err)
	}

	return &ThemeFile{
		path: filename,
		data: data,
	}
}

type ThemeFile struct {
	path string
	data []byte
}

func (themeFile *ThemeFile) StoreOnDisc(baseFolder string) (success bool, err error) {
	if !fsutil.CreateDirectory(baseFolder) {
		return false, fmt.Errorf("Unable to create the base folder for the themes: %q", baseFolder)
	}

	filePath := filepath.Join(baseFolder, themeFile.path)
	directory := filepath.Dir(filePath)

	if !fsutil.CreateDirectory(directory) {
		return false, fmt.Errorf("Unable to create folder %q for theme file %q.", directory, filePath)
	}

	if err := ioutil.WriteFile(filePath, themeFile.data, 0600); err != nil {
		return false, fmt.Errorf("Unable to create theme file %q in folder %q.", themeFile.path, baseFolder)
	}

	return true, nil
}
