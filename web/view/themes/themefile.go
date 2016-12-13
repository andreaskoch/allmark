// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package themes

import (
	"github.com/andreaskoch/allmark/common/util/fsutil"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func newFileFromText(uri, text string) *ThemeFile {
	return &ThemeFile{
		path: uri,
		data: []byte(text),
	}
}

func newFileFromBase64(uri, base64Text string) *ThemeFile {

	data, err := base64.StdEncoding.DecodeString(base64Text)
	if err != nil {
		panic(err)
	}

	return &ThemeFile{
		path: uri,
		data: data,
	}
}

type ThemeFile struct {
	path string
	data []byte
}

// Get the path of the theme file (e.g. "favicon.ico").
func (file *ThemeFile) Path() string {
	return file.path
}

// Get the data of the theme file.
func (file *ThemeFile) Data() []byte {
	return file.data
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
