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
	pathSeperator := string(FilesystemPathSeperator)
	relativePath := strings.Replace(pather.Path(), provider.basePath, "", 1)
	relativePath = pathSeperator + strings.TrimLeft(relativePath, pathSeperator)
	return relativePath
}
