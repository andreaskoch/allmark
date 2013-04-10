package html

import (
	"fmt"
	pathpackage "github.com/andreaskoch/allmark/path"
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

func renderImageGalleries(item *repository.Item) *repository.Item {

	for lineNumber, line := range item.RawLines {
		if galleryFound, transformedLine := renderImageGallery(item, line); galleryFound {
			item.RawLines[lineNumber] = transformedLine
		}
	}

	return item
}

func renderImageGallery(item *repository.Item, line string) (imageGalleryFound bool, transformedLine string) {
	found, matches := util.IsMatch(line, imageGalleryPattern)
	if !found || (found && len(matches) != 3) {
		return false, line
	}

	// parameters
	originalText := strings.TrimSpace(matches[0])
	galleryTitle := strings.TrimSpace(matches[1])
	path := strings.TrimSpace(matches[2])

	// create image gallery code
	files := item.Files.GetFilesByPath(path, isImageFile)

	imageLinks := getImageLinks(galleryTitle, item, files)
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
	return true, strings.Replace(line, originalText, imageGalleryCode, 1)
}

func getImageLinks(galleryTitle string, item *repository.Item, files []*repository.File) []string {

	pathProvider := pathpackage.NewProvider(item.Directory())
	numberOfFiles := len(files)
	imagelinks := make([]string, numberOfFiles, numberOfFiles)

	for index, file := range files {

		imagePath := pathProvider.GetFileRoute(file)
		imageTitle := fmt.Sprintf("%s - %s (Image %v of %v)", galleryTitle, getFileTitle(file), index+1, numberOfFiles)

		imagelinks[index] = fmt.Sprintf(`<a href="%s" title="%s"><img src="%s" /></a>`, imagePath, imageTitle, imagePath)
	}

	return imagelinks
}

func getFileTitle(pather pathpackage.Pather) string {
	fileName := filepath.Base(pather.Path())
	fileExtension := filepath.Ext(pather.Path())

	// remove the file extension from the file name
	filenameWithoutExtension := fileName[0:(strings.LastIndex(fileName, fileExtension))]

	return filenameWithoutExtension
}

func isImageFile(pather pathpackage.Pather) bool {
	fileExtension := strings.ToLower(filepath.Ext(pather.Path()))
	switch fileExtension {
	case ".png", ".gif", ".jpeg", ".jpg", ".svg", ".tiff":
		return true
	default:
		return false
	}

	panic("Unreachable")
}
