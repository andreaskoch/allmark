// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package initialization

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/web/view/templates"
	"github.com/andreaskoch/allmark2/web/view/themes"
)

func Initialize(baseFolder string) (success bool, err error) {
	config := config.Get(baseFolder)

	// create config
	if _, err := config.Save(); err != nil {
		return false, fmt.Errorf("Error while creating configuration file %q. Error: ", config.Filepath(), err)
	}

	fmt.Printf("Configuration file created at %q.\n", config.Filepath())

	// create theme
	themeFolder := config.ThemeFolder()
	if success, err := createTheme(themeFolder); !success {
		return false, fmt.Errorf("%s", err)
	}

	fmt.Printf("Theme stored in folder %q.\n", themeFolder)

	// create templates
	templateFolder := config.TemplatesFolder()
	if success, err := createTemplates(templateFolder); !success {
		return false, fmt.Errorf("%s", err)
	}

	fmt.Printf("Templates stored in folder %q.\n", templateFolder)
	return true, nil
}

func createTheme(baseFolder string) (success bool, err error) {
	return themes.GetTheme().StoreOnDisc(baseFolder)
}

func createTemplates(baseFolder string) (success bool, err error) {
	templateProvider := templates.NewProvider(baseFolder)
	return templateProvider.StoreTemplatesOnDisc()
}
