// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

import (
	"fmt"
	"github.com/andreaskoch/allmark/util"
	"github.com/andreaskoch/allmark/watcher"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	TemplateFileExtension = ".gohtml"
)

func NewTemplate(templateFolder string, name string, text string) *Template {

	// assemble the file path
	templateFilename := name + TemplateFileExtension
	templateFilePath := filepath.Join(templateFolder, templateFilename)

	// create a file change handler
	changeHandler, err := watcher.NewChangeHandler(templateFilePath)
	if err != nil {
		panic(fmt.Sprintf("Could not create a change handler for template %q.\nError: %s\n", templateFilePath, err))
	}

	return &Template{
		ChangeHandler: changeHandler,

		name: name,
		text: text,
		path: templateFilePath,
	}
}

type Template struct {
	*watcher.ChangeHandler

	name string
	text string
	path string
}

func (template *Template) Name() string {
	return template.name
}

func (template *Template) Text() string {

	if !util.FileExists(template.path) {
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
	if success, _ := util.CreateFile(path); !success {
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
