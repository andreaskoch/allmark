// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"fmt"
	"github.com/andreaskoch/allmark/util"
	"os/user"
)

type Config struct {
}

func GetConfig(repositoryPath string) Config {

}

// Get the current users home directory path
func getUserHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", errors.New(fmt.Sprintf("Cannot determine the current users home direcotry. Error: %s", err))
	}

	return usr.HomeDir, nil
}
