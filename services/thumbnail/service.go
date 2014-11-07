// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package thumbnail

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/dataaccess"
	"github.com/andreaskoch/allmark2/services/imageconversion"
	"io"
	"path/filepath"
	"strings"
	"time"
)

func NewConversionService(logger logger.Logger, config config.Config, repository dataaccess.Repository) *ConversionService {

	// assemble the index file path
	indexFilePath := filepath.Join(config.MetaDataFolder(), "thumbnails")
	index, err := loadIndex(indexFilePath)
	if err != nil {
		logger.Debug("No thumbnail index loaded (%s). Creating a new one.", err.Error())
	}

	// prepare the target folder
	targetFolder := filepath.Join(config.MetaDataFolder(), "thumbnails")
	logger.Debug("Creating a thumbnail folder at %q.", targetFolder)
	if !fsutil.CreateDirectory(targetFolder) {
		logger.Warn("Could not create the thumbnail folder %q", targetFolder)
		return nil
	}

	// create a new conversion service
	conversionService := &ConversionService{
		logger:     logger,
		config:     config,
		repository: repository,

		// thumbnail index
		indexFilePath: indexFilePath,
		index:         index,

		thumbnailFolder: targetFolder,
	}

	// start the conversion
	go conversionService.startConversion()

	return conversionService
}

type ConversionService struct {
	logger     logger.Logger
	config     config.Config
	repository dataaccess.Repository

	indexFilePath string
	index         Index

	thumbnailFolder string
}

func (conversion *ConversionService) startConversion() {

	conversion.createThumbnails()

	updateChannel := make(chan bool, 1)
	conversion.repository.AfterReindex(updateChannel)

	// refresh control
	go func() {
		for {
			select {
			case <-updateChannel:
				conversion.logger.Debug("Refreshing thumbnails")
				conversion.createThumbnails()
			}
		}
	}()

}

func (conversion *ConversionService) createThumbnails() {
ItemLoop:
	for _, item := range conversion.repository.Items() {

		for _, file := range item.Files() {
			conversion.createThumbnail(file)
			time.Sleep(5 * time.Second)
			break ItemLoop
		}
	}
}

func (conversion *ConversionService) createThumbnail(file *dataaccess.File) {

	// get the mime type
	mimeType, err := file.MimeType()
	if err != nil {
		conversion.logger.Warn("Unable to detect mime type for file. Error: %s", err.Error())
		return
	}

	// check type
	if !strings.HasPrefix(mimeType, "image/") {
		conversion.logger.Warn("%q is not an image", file)
		return
	}

	maxWidth := 100
	maxHeight := 100

	filename := fmt.Sprintf("%s-%v-%v", file.Name(), maxWidth, maxHeight)
	filePath := filepath.Join(conversion.thumbnailFolder, filename)
	target, err := fsutil.OpenFile(filePath)
	if err != nil {
		conversion.logger.Warn("Unable to detect mime type for file. Error: %s", err.Error())
		return
	}

	// convert the image
	conversionError := file.Data(func(content io.ReadSeeker) error {
		return imageconversion.Thumb(content, mimeType, 100, 100, target)
	})

	if conversionError != nil {
		conversion.logger.Warn("Unable to create thumbnail for file %q. Error: %s", file, err.Error())
		return
	}
}
