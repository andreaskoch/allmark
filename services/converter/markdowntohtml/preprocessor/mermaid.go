// Copyright 2020 Thomas Marschall. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package preprocessor

import (
	"github.com/elWyatt/allmark/common/paths"
	"github.com/elWyatt/allmark/model"
	"fmt"
	"regexp"
	"strings"
)

var (
	// ```mermaid<multiline code>```
	mermaidMarkdownExtensionPattern = regexp.MustCompile(`(?s)\x60{3}mermaid(.*?)\x60{3}`)
)

func newMermaidExtension(pathProvider paths.Pather, files []*model.File) *MermaidTableExtension {
	return &MermaidTableExtension{
		pathProvider: pathProvider,
		files:        files,
	}
}

type MermaidTableExtension struct {
	pathProvider paths.Pather
	files        []*model.File
}

func (converter *MermaidTableExtension) Convert(markdown string) (convertedContent string, converterError error) {

	convertedContent = markdown

	for _, match := range mermaidMarkdownExtensionPattern.FindAllStringSubmatch(convertedContent, -1) {

		if len(match) != 2 {
			continue
		}

		// parameters
		originalText := strings.TrimSpace(match[0])
		code := strings.TrimSpace(match[1])

		// get the code
		renderedCode := converter.getMermaidDiv(code)

		// replace markdown
		convertedContent = strings.Replace(convertedContent, originalText, renderedCode, 1)

	}

	return convertedContent, nil
}

func (converter *MermaidTableExtension) getMermaidDiv(code string) string {

	// div container with mermaid code
	divCode := fmt.Sprintf(`<div class="mermaid">%s</div>`, code)
	return divCode

}
