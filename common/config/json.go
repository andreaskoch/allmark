// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"encoding/json"
	"io"
)

type JSONSerializer struct{}

func NewJSONSerializer() *JSONSerializer {
	return &JSONSerializer{}
}

func (JSONSerializer) SerializeConfig(writer io.Writer, config *Config) error {
	bytes, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}

func (JSONSerializer) DeserializeConfig(reader io.Reader) (*Config, error) {
	decoder := json.NewDecoder(reader)
	var config *Config
	err := decoder.Decode(&config)
	return config, err
}
