// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package initialization

import (
	"allmark.io/modules/common/config"
	"allmark.io/modules/common/util/fsutil"
	"allmark.io/modules/web/view/templates"
	"allmark.io/modules/web/view/themes"
	"fmt"
)

func Initialize(baseFolder string) (success bool, err error) {
	config := config.Get(baseFolder)

	// create config
	if _, err := config.Save(); err != nil {
		return false, fmt.Errorf("Error while creating configuration file %q. Error: %s", config.Filepath(), err.Error())
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

	// empty basic-authentication file
	htpasswdFile := config.AuthenticationFilePath()
	if !fsutil.FileExists(htpasswdFile) {
		if _, err := fsutil.CreateFile(htpasswdFile); err != nil {
			return false, fmt.Errorf("Could not create a authentication user store. Error: %s", err.Error())
		}

		fmt.Printf("Created an empty authentication user store file: %q\n", config.AuthenticationFilePath())
	}

	// certs directory
	certificateDirectory := config.CertificateDirectory()
	if !fsutil.DirectoryExists(certificateDirectory) {
		if created := fsutil.CreateDirectory(certificateDirectory); !created {
			return false, fmt.Errorf("Could not create the certifcates directory: %q", config.CertificateDirectory())
		}

		fmt.Printf("Created the certifcates directory: %q\n", config.CertificateDirectory())
	}

	// ssl-certificates
	certFilePath, keyFilePath, created := config.CertificateFilePaths()
	if created {
		fmt.Printf("Created a certificate (%s, %s)\n", certFilePath, keyFilePath)
	}

	return true, nil
}

func createTheme(baseFolder string) (success bool, err error) {
	return themes.GetTheme().StoreOnDisc(baseFolder)
}

func createTemplates(baseFolder string) (success bool, err error) {
	templateProvider := templates.NewProvider(baseFolder)
	return templateProvider.StoreTemplatesOnDisc()
}
