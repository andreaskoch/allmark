// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	TemplateFileExtension = ".gohtml"
)

func NewTemplate(templateFolder, name, text string) *Template {

	// assemble the file path
	templateFilename := name + TemplateFileExtension
	templateFilePath := filepath.Join(templateFolder, templateFilename)

	// create a new template
	template := &Template{
		name: name,
		text: text,
		path: templateFilePath,
	}

	return template
}

type Template struct {
	name string
	text string
	path string
}

func (template *Template) Name() string {
	return template.name
}

func (template *Template) Text() string {

	if !fsutil.FileExists(template.path) {
		return template.text
	}

	file, err := os.Open(template.path)
	if err != nil {
		fmt.Printf("Could not open the template file %q.", template.path)
		return template.text
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("Could not read the template file %q.", template.path)
		return template.text
	}

	return string(bytes)
}

func (template *Template) StoreOnDisc() (success bool, err error) {

	path := template.path

	// make sure the directory exists
	if success, _ := fsutil.CreateFile(path); !success {
		return false, fmt.Errorf("Could not create the file %q.", path)
	}

	file, err := os.OpenFile(path, os.O_WRONLY, 0600)
	if err != nil {
		return false, err
	}

	defer file.Close()

	if _, err := file.WriteString(template.text); err != nil {
		return false, err
	}

	return true, nil
}
