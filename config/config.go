// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"bufio"
	"fmt"
	"github.com/andreaskoch/allmark/util"
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
	ThemeFolderName string
	Http            Http
}

type Config struct {
	Server Server

	metaDataFolder string
}

func (config *Config) SetMetaDataFolder(folder string) *Config {
	config.metaDataFolder = filepath.Join(folder, MetaDataFolderName)
	return config
}

func (config *Config) MetaDataFolder() string {
	return config.metaDataFolder
}

func (config *Config) ThemeFolder() string {
	return filepath.Join(config.MetaDataFolder(), config.Server.ThemeFolderName)
}

func emptyConfig() *Config {
	return &Config{}
}

func Initialize(repositoryPath string) {

	config := GetConfig(repositoryPath)

	// create config
	configurationFilePath := getConfigurationFilePath(repositoryPath)
	if ok, err := writeConfigToFile(config, configurationFilePath); !ok {
		fmt.Fprintf(os.Stderr, "Error while creating configuration file %q. Error: ", configurationFilePath, err)
	}

	// create theme
	themeFolder := config.ThemeFolder()
	if !util.CreateDirectory(themeFolder) {
		fmt.Fprintf(os.Stderr, "Unable to create theme folder %q.", themeFolder)
	}
}

func GetConfig(repositoryPath string) *Config {

	// return the local config
	localConfigurationFile := getConfigurationFilePath(repositoryPath)
	if localConfig, err := readConfigFromFile(localConfigurationFile); err == nil {
		return localConfig.SetMetaDataFolder(repositoryPath)
	}

	// return the global config
	if homeDirectory, homeDirError := getUserHomeDir(); homeDirError == nil {
		globalConfigurationFile := getConfigurationFilePath(homeDirectory)
		if globalConfig, configError := readConfigFromFile(globalConfigurationFile); configError == nil {
			return globalConfig.SetMetaDataFolder(homeDirectory)
		}
	}

	// return the default config
	defaultConfig := &Config{
		Server: Server{
			ThemeFolderName: ThemeFolderName,
			Http: Http{
				Port: 8080,
			},
		},
	}

	return defaultConfig.SetMetaDataFolder(repositoryPath)
}

func writeConfigToFile(config *Config, path string) (success bool, err error) {
	serializer := NewJSONSerializer()

	// create or overwrite the config file
	if created, err := util.CreateFile(path); !created {
		return false, fmt.Errorf("Could not create configuration file %q. Error: ", path, err)
	}

	// open the file for writing
	file, err := os.OpenFile(path, os.O_WRONLY, 0776)
	writer := bufio.NewWriter(file)

	defer func() {
		writer.Flush()
		file.Close()
	}()

	// serialize the config
	if serializationError := serializer.SerializeConfig(writer, config); serializationError != nil {
		return false, fmt.Errorf("Error while saving configuration %#v to file %q. Error: %v", config, path, serializationError)
	}

	return true, nil
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

// Get the current users home directory path
func getUserHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("Cannot determine the current users home direcotry. Error: %s", err)
	}

	return usr.HomeDir, nil
}
