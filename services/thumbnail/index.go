// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package thumbnail

import (
	"encoding/json"
	"fmt"
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

	// deserialize config
	serializer := NewJSONSerializer()
	index, err := serializer.DeserializeIndex(file)
	if err != nil {
		return Index{}, fmt.Errorf("Could not deserialize the index file %q. Error: %s", indexFilePath, err)
	}

	return index, nil
}

type JSONSerializer struct{}

func NewJSONSerializer() *JSONSerializer {
	return &JSONSerializer{}
}

func (JSONSerializer) SerializeIndex(writer io.Writer, index *Index) error {
	bytes, err := json.MarshalIndent(index, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}

func (JSONSerializer) DeserializeIndex(reader io.Reader) (Index, error) {
	decoder := json.NewDecoder(reader)
	var index Index
	err := decoder.Decode(index)
	return index, err
}

type Route string

type Thumb struct {
	MaxWidth  uint   `json:"maxWidth"`
	MaxHeight uint   `json:"maxHeight"`
	Path      string `json:"path"`
}

type Thumbs []Thumb

type Index map[Route]Thumbs
