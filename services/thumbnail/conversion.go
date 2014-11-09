// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package thumbnail

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/shutdown"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/dataaccess"
	"github.com/andreaskoch/allmark2/services/imageconversion"
	"io"
	"path/filepath"
	"time"
)

var (
	SizeSmall = ThumbDimension{
		MaxWidth:  200,
		MaxHeight: 0,
	}

	SizeMedium = ThumbDimension{
		MaxWidth:  400,
		MaxHeight: 0,
	}

	SizeLarge = ThumbDimension{
		MaxWidth:  800,
		MaxHeight: 0,
	}
)

func NewConversionService(logger logger.Logger, repository dataaccess.Repository, targetFolder string, thumbnailIndex *Index) *ConversionService {

	// create a new conversion service
	conversionService := &ConversionService{
		logger:     logger,
		repository: repository,

		isRunning:       true,
		conversionQueue: make(chan bool, 10),

		index:           thumbnailIndex,
		thumbnailFolder: targetFolder,
	}

	// start the conversion
	go conversionService.startConversion()

	// stop the conversion on shutdown
	shutdown.Register(func() error {
		logger.Info("Stopping the conversion process")
		conversionService.isRunning = false
		return nil
	})

	return conversionService
}

type ConversionService struct {
	logger     logger.Logger
	repository dataaccess.Repository

	isRunning       bool
	conversionQueue chan bool

	index           *Index
	thumbnailFolder string
}

func (conversion *ConversionService) startConversion() {

	updateChannel := make(chan bool, 1)
	conversion.repository.AfterReindex(updateChannel)

	// refresh control
	go func() {
		for conversion.isRunning {
			select {
			case <-updateChannel:
				conversion.conversionQueue <- true
			}
		}
	}()

	// thumbnail conversion
	go func() {

		for conversion.isRunning {

			select {
			case <-conversion.conversionQueue:

				for _, item := range conversion.repository.Items() {

					for _, file := range item.Files() {

						// create the thumbnail
						conversion.createThumbnail(file, SizeSmall)
						conversion.createThumbnail(file, SizeMedium)
						conversion.createThumbnail(file, SizeLarge)

						// wait before processing the next image
						time.Sleep(500 * time.Millisecond)
					}
				}

			}

		}

	}()

	// trigger the first conversion process
	conversion.conversionQueue <- true

}

func (conversion *ConversionService) createThumbnail(file *dataaccess.File, dimensions ThumbDimension) {

	// get the mime type
	mimeType, err := file.MimeType()
	if err != nil {
		conversion.logger.Warn("Unable to detect mime type for file. Error: %s", err.Error())
		return
	}

	// check the mime type
	if !imageconversion.MimeTypeIsSupported(mimeType) {
		conversion.logger.Debug("The mime-type %q is currently not supported.", mimeType)
		return
	}

	// determine the file name
	fileExtension := imageconversion.GetFileExtensionFromMimeType(mimeType)
	filename := fmt.Sprintf("%s-%v-%v.%s", file.Id(), dimensions.MaxWidth, dimensions.MaxHeight, fileExtension)

	// assemble the full file route
	fullFileRoute, err := route.Combine(file.Parent(), file.Route())
	if err != nil {
		conversion.logger.Warn("Unable to combine routes %q and %q.", file.Parent(), file.Route())
		return
	}

	thumb := newThumb(fullFileRoute, filename, dimensions)

	// check the index
	if conversion.isInIndex(thumb) {
		conversion.logger.Debug("Thumb %q already available in the index", thumb.String())
		return
	}

	// determine the file path
	filePath := filepath.Join(conversion.thumbnailFolder, filename)

	// open the target file
	target, fileError := fsutil.OpenFile(filePath)
	if fileError != nil {
		conversion.logger.Warn("Unable to detect mime type for file. Error: %s", fileError.Error())
		return
	}

	defer target.Close()

	// convert the image
	conversionError := file.Data(func(content io.ReadSeeker) error {
		return imageconversion.Resize(content, mimeType, dimensions.MaxWidth, dimensions.MaxHeight, target)
	})

	// handle errors
	if conversionError != nil {
		conversion.logger.Warn("Unable to create thumbnail for file %q. Error: %s", file, conversionError.Error())
		return
	}

	// add to index
	conversion.addToIndex(thumb)
	conversion.logger.Debug("Adding Thumb %q to index", thumb.String())
}

func (conversion *ConversionService) isInIndex(thumb Thumb) bool {

	// check if there are thumb for the route
	thumbs, entryExists := conversion.index.GetThumbs(thumb.Route)
	if !entryExists {
		return false
	}

	// check if there is a thumb with that dimensions
	if _, thumbExists := thumbs[thumb.Dimensions.String()]; thumbExists {
		// check if the file exists
		return fsutil.FileExists(filepath.Join(conversion.thumbnailFolder, thumb.Path))

	}

	return false
}

func (conversion *ConversionService) addToIndex(thumb Thumb) {
	thumbs, entryExists := conversion.index.GetThumbs(thumb.Route)
	if !entryExists {
		thumbs = make(Thumbs)
	}

	thumbs[thumb.Dimensions.String()] = thumb
	conversion.index.SetThumbs(thumb.Route, thumbs)
}
