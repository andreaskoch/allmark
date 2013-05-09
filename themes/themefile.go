// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package themes

import (
	"fmt"
	"github.com/andreaskoch/allmark/util"
	"os"
	"path/filepath"
)

type ThemeFile struct {
	Filename string
	Content  string
}

func (themeFile *ThemeFile) StoreOnDisc(baseFolder string) (success bool, err error) {
	if !util.CreateDirectory(baseFolder) {
		return false, fmt.Errorf("Unable to create theme folder %q.", baseFolder)
	}

	filePath := filepath.Join(baseFolder, themeFile.Filename)
	file, err := os.Create(filePath)
	if err != nil {
		return false, fmt.Errorf("Unable to create theme file %q in folder %q.", themeFile.Filename, baseFolder)
	}

	defer file.Close()

	file.WriteString(themeFile.Content)

	return true, nil
}
