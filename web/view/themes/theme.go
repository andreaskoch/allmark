// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package themes

func GetTheme() *Theme {
	return defaultTheme
}

type Theme struct {
	Name  string
	Files []*ThemeFile
}

// Get the theme file that matches the specified uri (e.g. "favicon.ico"); returns nil if the theme file was not found.
func (theme *Theme) Get(uri string) *ThemeFile {
	for _, themeFile := range theme.Files {
		if themeFile.Path() == uri {
			return themeFile
		}
	}

	return nil
}

func (theme *Theme) StoreOnDisc(baseFolder string) (success bool, err error) {
	for _, file := range theme.Files {
		if ok, err := file.StoreOnDisc(baseFolder); !ok {
			return false, err
		}
	}

	return true, nil
}
