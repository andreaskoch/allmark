// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package common

import (
	"time"
)

type LastModifiedProviderFunc func() (time.Time, error)

type ContentProviderFunc func() ([]byte, error)

type HashProviderFunc func() (string, error)
