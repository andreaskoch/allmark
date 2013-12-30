// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

import (
	"fmt"
	"github.com/andreaskoch/allmark/util"
	"github.com/andreaskoch/go-fswatch"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	TemplateFileExtension = ".gohtml"
)

func NewTemplate(templateFolder, name, text string, modified chan bool) *Template {

	// assemble the file path
	templateFilename := name + TemplateFileExtension
	templateFilePath := filepath.Join(templateFolder, templateFilename)

	// create a new template
	template := &Template{
		Modified: make(chan bool),
		Moved:    make(chan bool),

		name: name,
		text: text,
		path: templateFilePath,
	}

	// look for changes
	if util.FileExists(templateFilePath) {
		go func() {
			fileWatcher := fswatch.NewFileWatcher(template.path).Start()

			for fileWatcher.IsRunning() {

				select {
				case <-fileWatcher.Modified:

					fmt.Printf("Template %q changed.\n", templateFilePath)

					go func() {
						modified <- true
					}()

				case <-fileWatcher.Moved:

					fmt.Printf("Template %q moved.\n", templateFilePath)

					go func() {
						modified <- true
					}()
				}

			}
		}()
	}

	return template
}

type Template struct {
	Modified chan bool
	Moved    chan bool

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
