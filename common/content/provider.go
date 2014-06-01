// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package content

import (
	"io"
	"time"
)

type MimeTypeProviderFunc func() (string, error)

type LastModifiedProviderFunc func() (time.Time, error)

type DataProviderFunc func(contentReader func(content io.ReadSeeker) error) error

type HashProviderFunc func() (string, error)

type ChangeEventProviderFunc func() chan ChangeEvent

func NewProvider(mimeType MimeTypeProviderFunc, data DataProviderFunc, hash HashProviderFunc, lastModified LastModifiedProviderFunc, changes ChangeEventProviderFunc) *ContentProvider {
	return &ContentProvider{
		mimeTypeProviderFunc:     mimeType,
		dataProviderFunc:         data,
		hashProviderFunc:         hash,
		lastModifiedProviderFunc: lastModified,
		changeEventProviderFunc:  changes,
	}
}

type ContentProvider struct {
	mimeTypeProviderFunc     MimeTypeProviderFunc
	dataProviderFunc         DataProviderFunc
	hashProviderFunc         HashProviderFunc
	lastModifiedProviderFunc LastModifiedProviderFunc
	changeEventProviderFunc  ChangeEventProviderFunc
}

func (provider *ContentProvider) Data(contentReader func(content io.ReadSeeker) error) error {
	return provider.dataProviderFunc(contentReader)
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

func (provider *ContentProvider) ChangeEvent() chan ChangeEvent {
	return provider.changeEventProviderFunc()
}
