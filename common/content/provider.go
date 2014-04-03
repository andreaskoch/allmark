// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package content

import (
	"time"
)

func NewProvider(mimeType MimeTypeProviderFunc, data DataProviderFunc, hash HashProviderFunc, lastModified LastModifiedProviderFunc) *ContentProvider {
	return &ContentProvider{
		mimeTypeProviderFunc:     mimeType,
		dataProviderFunc:         data,
		hashProviderFunc:         hash,
		lastModifiedProviderFunc: lastModified,
	}
}

type MimeTypeProviderFunc func() (string, error)

type LastModifiedProviderFunc func() (time.Time, error)

type DataProviderFunc func() ([]byte, error)

type HashProviderFunc func() (string, error)

type ContentProvider struct {
	mimeTypeProviderFunc     MimeTypeProviderFunc
	dataProviderFunc         DataProviderFunc
	hashProviderFunc         HashProviderFunc
	lastModifiedProviderFunc LastModifiedProviderFunc
}

func (provider *ContentProvider) Data() ([]byte, error) {
	return provider.dataProviderFunc()
}

func (provider *ContentProvider) Hash() (string, error) {
	return provider.hashProviderFunc()
}

func (provider *ContentProvider) LastModified() (time.Time, error) {
	return provider.lastModifiedProviderFunc()
}

func (provider *ContentProvider) MimeType() (string, error) {
	return provider.mimeTypeProviderFunc()
}
