package renderer

import (
	"fmt"
	p "github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/util"
	"path/filepath"
	"strings"
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
	found, matches := util.IsMatch(line, ImageGalleryPattern)
	if !found || (found && len(matches) != 3) {
		return false, line
	}

	// parameters
	originalText := matches[0]
	descriptionText := matches[1]
	path := matches[2]

	// create image gallery code
	files := item.Files.GetFilesByPath(path, func(pather p.Pather) bool {
		fileExtension := strings.ToLower(filepath.Ext(pather.Path()))
		switch fileExtension {
		case ".png", ".gif", ".jpeg", ".jpg", ".svg", ".tiff":
			return true
		default:
			return false
		}

		panic("Unreachable")
	})

	imageLinks := getImageLinks(item, files)
	imageGalleryCode := fmt.Sprintf(`<div class="imagegallery">
				<header>
					<span>%s</span>
				</header>
				%s
			</div>`, descriptionText, strings.Join(imageLinks, "\n"))

	// replace markdown with image gallery
	return true, strings.Replace(line, originalText, imageGalleryCode, 1)
}

func getImageLinks(item *repository.Item, files []*repository.File) []string {
	pathProvider := p.NewProvider(item.Directory())
	imagelinks := make([]string, len(files), len(files))
	for index, file := range files {
		imagelinks[index] = getImageLink(pathProvider, file)
	}
	return imagelinks
}

func getImageLink(pathProvider *p.Provider, file *repository.File) string {
	fileRoute := pathProvider.GetFileRoute(file)
	return fmt.Sprintf(`<img src="%s" />`, fileRoute)
}
