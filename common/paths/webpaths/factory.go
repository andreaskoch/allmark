// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webpaths

import (
	"github.com/andreaskoch/allmark2/common/paths"
)

func NewFactory() *PatherFactory {
	return &PatherFactory{}
}

type PatherFactory struct {
}

func (factory *PatherFactory) Absolute(base string) paths.Pather {
	return newAbsoluteWebPathProvider(base)
}

func (factory *PatherFactory) Relative() paths.Pather {
	return newRelativeWebPathProvider()
}
