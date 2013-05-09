// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package themes

func GetTheme() *Theme {
	return &defaultTheme
}

type Theme struct {
	Files []*ThemeFile
}

func (theme *Theme) StoreOnDisc(baseFolder string) (success bool, err error) {
	for _, file := range theme.Files {
		if ok, err := file.StoreOnDisc(baseFolder); !ok {
			return false, err
		}
	}

	return true, nil
}
