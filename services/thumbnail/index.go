// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package thumbnail

import (
	"encoding/json"
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/shutdown"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"io"
	"os"
)

func NewIndex(logger logger.Logger, indexFilePath string) *Index {

	// assemble the index file path
	index, err := loadIndex(indexFilePath)
	if err != nil {
		logger.Debug("No thumbnail index loaded (%s). Creating a new one.", err.Error())
	}

	// save the index on shutdown
	shutdown.Register(func() error {
		logger.Info("Saving the index")
		return saveIndex(index, indexFilePath)
	})

	return index
}

func emptyIndex() *Index {
	return &Index{
		make(map[string]Thumbs),
	}
}

func loadIndex(indexFilePath string) (*Index, error) {

	if !fsutil.FileExists(indexFilePath) {
		return emptyIndex(), fmt.Errorf("The index file %q does not exist.", indexFilePath)
	}

	// check if file can be accessed
	file, err := os.Open(indexFilePath)
	if err != nil {
		return emptyIndex(), fmt.Errorf("Cannot read index file %q. Error: %s", indexFilePath, err)
	}

	defer file.Close()

	// deserialize the index
	serializer := newIndexSerializer()
	index, err := serializer.DeserializeIndex(file)
	if err != nil {
		return emptyIndex(), fmt.Errorf("Could not deserialize the index file %q. Error: %s", indexFilePath, err)
	}

	return index, nil
}

func saveIndex(index *Index, indexFilePath string) error {
	file, fileError := fsutil.OpenFile(indexFilePath)
	if fileError != nil {
		return fmt.Errorf("Cannot save index to file %q. Error: %s", indexFilePath, fileError.Error())
	}

	defer file.Close()

	// serialize the index
	serializer := newIndexSerializer()
	return serializer.SerializeIndex(file, index)
}

type indexSerializer struct{}

func newIndexSerializer() *indexSerializer {
	return &indexSerializer{}
}

func (indexSerializer) SerializeIndex(writer io.Writer, index *Index) error {
	bytes, err := json.MarshalIndent(index, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}

func (indexSerializer) DeserializeIndex(reader io.Reader) (*Index, error) {
	decoder := json.NewDecoder(reader)
	var index Index
	err := decoder.Decode(index)
	return &index, err
}

func newThumb(route route.Route, path string, maxWidth, maxHeight uint) Thumb {

	return Thumb{
		Route: route.Value(),
		Path:  path,
		Dimensions: ThumbDimension{
			MaxWidth:  maxWidth,
			MaxHeight: maxHeight,
		},
	}

}

type Thumb struct {
	Route      string         `json:"route"`
	Path       string         `json:"path"`
	Dimensions ThumbDimension `json:"dimensions"`
}

func (t Thumb) String() string {
	return fmt.Sprintf("%s (%s)", t.Path, t.Dimensions.String())
}

func (t Thumb) ThumbRoute() route.Route {
	thumbRoute, err := route.NewFromRequest(fmt.Sprintf("%s-%s", t.Route, t.Dimensions.String()))
	if err != nil {
		panic(err)
	}

	return thumbRoute
}

type ThumbDimension struct {
	MaxWidth  uint `json:"maxWidth"`
	MaxHeight uint `json:"maxHeight"`
}

func (t ThumbDimension) String() string {
	return fmt.Sprintf("maxWidth:%v-maxHeight:%v", t.MaxWidth, t.MaxHeight)
}

type Thumbs map[string]Thumb

func (thumbs Thumbs) GetThumbBySize(maxWidth, maxHeight uint) (Thumb, bool) {

	dimension := ThumbDimension{
		MaxWidth:  maxWidth,
		MaxHeight: maxHeight,
	}

	thumb, exists := thumbs[dimension.String()]
	return thumb, exists
}

type Index struct {
	Thumbs map[string]Thumbs `json:"thumbs"`
}

func (i *Index) GetThumbs(key string) (thumbs Thumbs, exists bool) {
	thumbs, exists = i.Thumbs[key]
	return thumbs, exists
}

func (i *Index) SetThumbs(key string, thumbs Thumbs) {
	i.Thumbs[key] = thumbs
}
