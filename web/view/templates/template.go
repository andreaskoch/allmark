// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templates

import (
	"github.com/andreaskoch/allmark/common/util/fsutil"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	TemplateFileExtension = ".gohtml"
)

// newTemplateDefinition creates a new template definition with the given parameters.
func newTemplateDefinition(templateFolder, name, text string) *templateDefinition {

	// assemble the file path
	templateFilename := name + TemplateFileExtension
	templateFilePath := filepath.Join(templateFolder, templateFilename)

	// create a new template definition
	templateDefinition := &templateDefinition{
		name: name,
		text: text,
		path: templateFilePath,
	}

	return templateDefinition
}

// A templateDefinition contains template code, its name and the path on disc.
type templateDefinition struct {
	name string
	text string
	path string
}

// Name returns the template name
func (template *templateDefinition) Name() string {
	return template.name
}

// Text returns the template code. If the template was found on disc it will return the code from disc.
// Otherwise it will return the default template code.
func (template *templateDefinition) Text() string {

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

// StoreOnDisc stores the current template definition to it's target path on disc.
func (template *templateDefinition) StoreOnDisc() (success bool, err error) {

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
