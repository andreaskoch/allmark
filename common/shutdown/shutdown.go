// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shutdown

import (
	"fmt"
)

var (
	callbacks = make([]func() error, 0)
)

func Register(callback func() error) {
	callbacks = append(callbacks, callback)
}

func Shutdown() {

	for _, callback := range callbacks {
		err := callback()
		if err != nil {
			fmt.Println(err.Error())
		}
	}

}
