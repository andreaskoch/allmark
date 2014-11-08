// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package thumbnail

import (
	"encoding/json"
	"fmt"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"io"
	"os"
)

func loadIndex(indexFilePath string) (Index, error) {

	if !fsutil.FileExists(indexFilePath) {
		return Index{}, fmt.Errorf("The index file %q does not exist.", indexFilePath)
	}

	// check if file can be accessed
	file, err := os.Open(indexFilePath)
	if err != nil {
		return Index{}, fmt.Errorf("Cannot read index file %q. Error: %s", indexFilePath, err)
	}

	defer file.Close()

	// deserialize the index
	serializer := newIndexSerializer()
	index, err := serializer.DeserializeIndex(file)
	if err != nil {
		return Index{}, fmt.Errorf("Could not deserialize the index file %q. Error: %s", indexFilePath, err)
	}

	return index, nil
}

func saveIndex(index Index, indexFilePath string) error {
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

func (indexSerializer) SerializeIndex(writer io.Writer, index Index) error {
	bytes, err := json.MarshalIndent(index, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}

func (indexSerializer) DeserializeIndex(reader io.Reader) (Index, error) {
	decoder := json.NewDecoder(reader)
	var index Index
	err := decoder.Decode(index)
	return index, err
}

func newThumb(route route.Route, path string, maxWidth, maxHeight uint) Thumb {

	return Thumb{
		Path: path,
		Dimensions: ThumbDimension{
			MaxHeight: maxHeight,
			MaxWidth:  maxWidth,
		},
	}

}

type Thumb struct {
	Dimensions ThumbDimension `json:"dimensions"`
	Path       string         `json:"path"`
}

func (t Thumb) String() string {
	return fmt.Sprintf("%s (%s)", t.Path, t.Dimensions.String())
}

type ThumbDimension struct {
	MaxWidth  uint `json:"maxWidth"`
	MaxHeight uint `json:"maxHeight"`
}

func (t ThumbDimension) String() string {
	return fmt.Sprintf("%v x %v", t.MaxWidth, t.MaxHeight)
}

type Thumbs map[string]Thumb

type Index map[string]Thumbs
