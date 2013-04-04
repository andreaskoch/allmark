package path

import (
	"github.com/andreaskoch/docs/repository"
	"os"
	"strings"
)

const (
	FilesystemPathSeperator = os.PathSeparator
	WebPathSeperator        = "/"
)

func NewProvider(basePath string) *Provider {
	return &Provider{
		basePath: basePath,
	}
}

type Provider struct {
	basePath string
}

func (provider *Provider) GetWebRoute(pather repository.Pather) string {
	filesystemSeperatorString := string(FilesystemPathSeperator)

	relativeFilepath := strings.Replace(pather.Path(), provider.basePath, "", 1)
	relativeFilepath = filesystemSeperatorString + strings.TrimLeft(relativeFilepath, filesystemSeperatorString)

	webRoute := strings.Replace(relativeFilepath, filesystemSeperatorString, "/", -1)
	return webRoute
}
