// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"fmt"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/util"
	"regexp"
	"strings"
)

var (
	// imagegallery: [*description text*](*folder path*)
	imageGalleryPattern = regexp.MustCompile(`imagegallery: \[([^\]]+)\]\(([^)]+)\)`)
)

func renderImageGalleries(fileIndex *repository.FileIndex, pathProvider *path.Provider, markdown string) string {

	for {

		found, matches := util.IsMatch(markdown, imageGalleryPattern)
		if !found || (found && len(matches) != 3) {
			break
		}

		// parameters
		originalText := strings.TrimSpace(matches[0])
		galleryTitle := strings.TrimSpace(matches[1])
		path := strings.TrimSpace(matches[2])

		// create image gallery code
		files := fileIndex.FilesByPath(path, isImageFile)

		imageLinks := getImageLinks(galleryTitle, files, pathProvider)
		imageGalleryCode := fmt.Sprintf(`<section class="imagegallery">
				<h1>%s</h1>
				<ol>
					<li>
					%s
					</li>
				</ol>
			</section>`, galleryTitle, strings.Join(imageLinks, "\n</li>\n<li>\n"))

		// replace markdown with image gallery
		markdown = strings.Replace(markdown, originalText, imageGalleryCode, 1)

	}

	return markdown
}

func getImageLinks(galleryTitle string, files []*repository.File, pathProvider *path.Provider) []string {

	numberOfFiles := len(files)
	imagelinks := make([]string, numberOfFiles, numberOfFiles)

	for index, file := range files {

		imagePath := pathProvider.GetWebRoute(file)
		imageTitle := fmt.Sprintf("%s - %s (Image %v of %v)", galleryTitle, getFileTitle(file), index+1, numberOfFiles)

		imagelinks[index] = fmt.Sprintf(`<a href="%s" title="%s"><img src="%s" /></a>`, imagePath, imageTitle, imagePath)
	}

	return imagelinks
}
