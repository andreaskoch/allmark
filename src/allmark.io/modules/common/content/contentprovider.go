// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package content

import (
	"io"
	"time"
)

type ContentProviderInterface interface {
	Data(contentReader func(content io.ReadSeeker) error) error
	Hash() (string, error)
	LastModified() (time.Time, error)
	MimeType() (string, error)
}

type MimeTypeProviderFunc func() (string, error)

type LastModifiedProviderFunc func() (time.Time, error)

type DataProviderFunc func(contentReader func(content io.ReadSeeker) error) error

type HashProviderFunc func() (string, error)

func NewContentProvider(mimeType MimeTypeProviderFunc, data DataProviderFunc, hash HashProviderFunc, lastModified LastModifiedProviderFunc) *ContentProvider {

	currentHash, err := hash()
	if err != nil {
		panic(err)
	}

	return &ContentProvider{
		mimeTypeProviderFunc:     mimeType,
		dataProviderFunc:         data,
		hashProviderFunc:         hash,
		lastModifiedProviderFunc: lastModified,
		lastHash:                 currentHash,
	}
}

type ContentProvider struct {
	mimeTypeProviderFunc     MimeTypeProviderFunc
	dataProviderFunc         DataProviderFunc
	hashProviderFunc         HashProviderFunc
	lastModifiedProviderFunc LastModifiedProviderFunc
	lastHash                 string
}

func (provider *ContentProvider) Data(contentReader func(content io.ReadSeeker) error) error {
	return provider.dataProviderFunc(contentReader)
}

func (provider *ContentProvider) Hash() (string, error) {
	hash, err := provider.hashProviderFunc()
	if err != nil {
		return hash, err
	}

	provider.lastHash = hash
	return hash, nil
}

func (provider *ContentProvider) LastModified() (time.Time, error) {
	return provider.lastModifiedProviderFunc()
}

func (provider *ContentProvider) MimeType() (string, error) {
	return provider.mimeTypeProviderFunc()
}

func (provider *ContentProvider) LastHash() string {
	return provider.lastHash
}
