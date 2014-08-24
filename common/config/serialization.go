// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"io"
)

type ConfigSerializer interface {
	SerializeConfig(writer io.Writer, config *Config) error
}

type ConfigDeserializer interface {
	DeserializeConfig(reader io.Reader) (*Config, error)
}
