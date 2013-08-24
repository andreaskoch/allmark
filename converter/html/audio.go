// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"fmt"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/util"
	"mime"
	"regexp"
	"strings"
)

var (
	// audio: [*description text*](*a link to an audio file*)
	audioPattern = regexp.MustCompile(`audio: \[([^\]]+)\]\(([^)]+)\)`)
)

func renderAudio(item *repository.Item, rawContent string) string {
	return convertAudioMarkdownExtension(rawContent, item.Files, item.FilePathProvider())
}

func convertAudioMarkdownExtension(markdown string, fileIndex *repository.FileIndex, pathProvider *path.Provider) string {

	for {

		found, matches := util.IsMatch(markdown, audioPattern)
		if !found || (found && len(matches) != 3) {
			break
		}

		// parameters
		originalText := strings.TrimSpace(matches[0])
		title := strings.TrimSpace(matches[1])
		path := strings.TrimSpace(matches[2])

		// get a renderer
		renderer := getAudioRenderer(title, path, fileIndex, pathProvider)

		// execute the renderer
		renderedCode := renderer()

		// replace markdown with link list
		markdown = strings.Replace(markdown, originalText, renderedCode, 1)

	}

	return markdown
}

func getAudioRenderer(title, path string, fileIndex *repository.FileIndex, pathProvider *path.Provider) func() string {

	// html5 audio file
	if isAudioFile, mimeType := isAudioFileLink(path); isAudioFile {
		return func() string {
			return renderAudioFileLink(title, path, mimeType)
		}
	}

	// return the fallback handler
	return func() string {
		return fmt.Sprintf(`<a href="%s" target="_blank" title="%s">%s</a>`, path, title, title)
	}
}

func isAudioFileLink(link string) (isAudioFile bool, mimeType string) {
	normalizedLink := strings.ToLower(link)
	fileExtension := normalizedLink[strings.LastIndex(normalizedLink, "."):]
	mimeType = mime.TypeByExtension(fileExtension)

	switch fileExtension {
	case ".mp3", ".ogg":
		return true, mimeType
	default:
		return false, ""
	}

	panic("Unreachable")
}

func renderAudioFileLink(title, link, mimetype string) string {
	return fmt.Sprintf(`<section class="audio audio-file">
		<h1><a href="%s" target="_blank" title="%s">%s</a></h1>
		<audio controls>
			<source src="%s" type="%s">
		</audio>
	</section>`, link, title, title, link, mimetype)
}
