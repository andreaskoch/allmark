// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package content

import (
	"fmt"
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

// NewContentProvider creates a new content provider with the given mimeType, data provider, hash provider and last modified provider.
func NewContentProvider(mimeType MimeTypeProviderFunc, data DataProviderFunc, hash HashProviderFunc, lastModified LastModifiedProviderFunc) (*ContentProvider, error) {

	currentHash, err := hash()
	if err != nil {
		return nil, fmt.Errorf("Cannot create content provider because hash() returned an error: %s", err.Error())
	}

	return &ContentProvider{
		mimeTypeProviderFunc:     mimeType,
		dataProviderFunc:         data,
		hashProviderFunc:         hash,
		lastModifiedProviderFunc: lastModified,
		lastHash:                 currentHash,
	}, nil
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
