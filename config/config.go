// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

const (
	MetaDataFolderName    = ".allmark"
	ConfigurationFileName = "config"
	ThemeFolderName       = "theme"
)

type Http struct {
	Port int
}

type Server struct {
	ThemeFolder string
	Http        Http
}

type Config struct {
	Server Server
}

func emptyConfig() *Config {
	return &Config{}
}

func Init(repositoryPath string) {
	config := GetConfig(repositoryPath)

	configurationFilePath := getConfigurationFilePath(repositoryPath)
	configurationFileDirectory := filepath.Dir(configurationFilePath)

	os.MkdirAll(configurationFileDirectory, os.ModeDir)

	writeConfigToFile(config, configurationFilePath)
}

func GetConfig(repositoryPath string) *Config {

	// return the local config
	localConfigurationFile := getConfigurationFilePath(repositoryPath)
	if localConfig, err := readConfigFromFile(localConfigurationFile); err == nil {
		return localConfig
	}

	// return the global config
	if homeDirectory, homeDirError := getUserHomeDir(); homeDirError == nil {
		globalConfigurationFile := getConfigurationFilePath(homeDirectory)
		if globalConfig, configError := readConfigFromFile(globalConfigurationFile); configError == nil {
			return globalConfig
		}
	}

	// return the default config
	return &Config{
		Server: Server{
			ThemeFolder: getThemeFolderPath(repositoryPath),
			Http: Http{
				Port: 8080,
			},
		},
	}
}

func writeConfigToFile(config *Config, path string) {
	serializer := NewJSONSerializer()

	file, err := os.OpenFile(path, os.O_WRONLY, 0776)
	if err != nil {
		panic(fmt.Sprintf("Could not create configuration file %q. Error: ", path, err))
	}

	writer := bufio.NewWriter(file)
	defer file.Close()

	if serializationError := serializer.SerializeConfig(writer, config); serializationError != nil {
		fmt.Printf("Error while saving configuration %#v to file %q. Error: %v", config, path, serializationError)
	}
}

func readConfigFromFile(path string) (*Config, error) {
	fileInfo, err := os.Open(path)
	if err != nil {
		return emptyConfig(), fmt.Errorf("Cannot read config file %q. Error: %s", path, err)
	}

	serializer := NewJSONSerializer()
	config, err := serializer.DeserializeConfig(fileInfo)
	if err != nil {
		return emptyConfig(), fmt.Errorf("Could not deserialize the configuration file %q. Error: %s", path, err)
	}

	return config, nil
}

func getConfigurationFilePath(folder string) string {
	return filepath.Join(folder, MetaDataFolderName, ConfigurationFileName)
}

func getThemeFolderPath(folder string) string {
	return filepath.Join(folder, MetaDataFolderName, ThemeFolderName)
}

// Get the current users home directory path
func getUserHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("Cannot determine the current users home direcotry. Error: %s", err)
	}

	return usr.HomeDir, nil
}
