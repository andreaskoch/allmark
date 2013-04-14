package html

import (
	"fmt"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/util"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// !imagegallery[*descriptionText*](*folderPath*)
	imageGalleryPattern = regexp.MustCompile(`!imagegallery\[([^\]]+)\]\(([^)]+)\)`)
)

func NewImageGalleryRenderer(markdown string, fileIndex *repository.FileIndex, pathProvider *path.Provider) func(text string) string {
	return func(text string) string {
		return renderImageGallery(markdown, fileIndex, pathProvider)
	}
}

func renderImageGallery(markdown string, fileIndex *repository.FileIndex, pathProvider *path.Provider) string {

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
		files := fileIndex.GetFilesByPath(path, isImageFile)

		imageLinks := getImageLinks(galleryTitle, files, pathProvider)
		imageGalleryCode := fmt.Sprintf(`<div class="imagegallery">
				<header>
					<span>%s</span>
				</header>
				<section>
				<ol>
					<li>
					%s
					</li>
				</ol>
				</section>
			</div>`, galleryTitle, strings.Join(imageLinks, "\n</li>\n<li>\n"))

		// replace markdown with image gallery
		markdown = strings.Replace(markdown, originalText, imageGalleryCode, 1)

	}

	return markdown
}

func getImageLinks(galleryTitle string, files []*repository.File, pathProvider *path.Provider) []string {

	numberOfFiles := len(files)
	imagelinks := make([]string, numberOfFiles, numberOfFiles)

	for index, file := range files {

		imagePath := pathProvider.GetFileRoute(file)
		imageTitle := fmt.Sprintf("%s - %s (Image %v of %v)", galleryTitle, getFileTitle(file), index+1, numberOfFiles)

		imagelinks[index] = fmt.Sprintf(`<a href="%s" title="%s"><img src="%s" /></a>`, imagePath, imageTitle, imagePath)
	}

	return imagelinks
}

func getFileTitle(pather path.Pather) string {
	fileName := filepath.Base(pather.Path())
	fileExtension := filepath.Ext(pather.Path())

	// remove the file extension from the file name
	filenameWithoutExtension := fileName[0:(strings.LastIndex(fileName, fileExtension))]

	return filenameWithoutExtension
}

func isImageFile(pather path.Pather) bool {
	fileExtension := strings.ToLower(filepath.Ext(pather.Path()))
	switch fileExtension {
	case ".png", ".gif", ".jpeg", ".jpg", ".svg", ".tiff":
		return true
	default:
		return false
	}

	panic("Unreachable")
}
