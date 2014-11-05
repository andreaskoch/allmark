// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package imageconversion

import (
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

func Thumb(source io.Reader, mimeType string, maxWidth, maxHeight int, target io.Writer) error {

	// read the source image
	img, err := decode(source, mimeType)
	if err != nil {
		return err
	}

	// resize the source image
	thumb := resize.Thumb(maxWidth, maxHeight, img, resize.Lanczos3)

	// write the thumbnail to the target
	return encode(mimeType, thumb, target)
}

func Resize(source io.Reader, mimeType string, width, height int, target io.Writer) error {

	// read the source image
	img, err := decode(source, mimeType)
	if err != nil {
		return err
	}

	// resize the source image
	thumb := resize.Resize(width, height, img, resize.Lanczos3)

	// write the thumbnail to the target
	return encode(mimeType, thumb, target)
}

func encode(mimeType string, thumb image.Image, target io.Writer) error {

	switch mimeType {

	case "image/png":
		return png.Encode(target, thumb)

	case "image/jpeg":
		return jpeg.Encode(target, thumb)

	case "image/gif":
		return gif.Encode(target, thumb)

	default:
		return nil, fmt.Errorf("Unsupported mime type %s", mimeType)

	}

	panic("Unreachable")
}

func decode(source io.Reader, mimeType string) (image.Image, error) {

	switch mimeType {

	case "image/png":
		return png.Decode(source)

	case "image/jpeg":
		return jpeg.Decode(source)

	case "image/gif":
		return gif.Decode(source)

	default:
		return nil, fmt.Errorf("Unsupported mime type %s", mimeType)

	}

	panic("Unreachable")
}
