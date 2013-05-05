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
	TemplatesFolderName   = "templates"
)

var isHomeDir func(directory string) bool

func init() {

	usr, err := user.Current()
	if err != nil {
		panic(fmt.Sprintf("Cannot determine the current users home direcotry. Error: %s", err))
	}

	isHomeDir = func(directory string) bool {
		return filepath.Clean(directory) == filepath.Clean(usr.HomeDir)
	}
}

func GetConfig(baseFolder string) *Config {

	// check if the base folder is the home dir
	if isHomeDir(baseFolder) {

		// return the global config
		if config, err := global(baseFolder).load(); err == nil {
			return config
		}

		return defaultConfig(baseFolder)
	}

	// return the local config
	if config, err := local(baseFolder).load(); err == nil {
		return config
	}

	// return the global config
	if config, err := global(baseFolder).load(); err == nil {
		return config
	}

	// return the default config
	return defaultConfig(baseFolder)
}

type Http struct {
	Port int
}

type Web struct {
	DefaultLanguage string
}

type Server struct {
	ThemeFolderName string
	Http            Http
}

type Config struct {
	Server Server
	Web    Web

	baseFolder      string
	metaDataFolder  string
	themeFolderBase string
	templatesFolder string
}

func (config *Config) BaseFolder() string {
	return config.baseFolder
}

func (config *Config) MetaDataFolder() string {
	return config.metaDataFolder
}

func (config *Config) TemplatesFolder() string {
	return config.templatesFolder
}

func (config *Config) Filepath() string {
	return filepath.Join(config.MetaDataFolder(), ConfigurationFileName)
}

func (config *Config) ThemeFolder() string {
	themeFolderName := ThemeFolderName
	if config.Server.ThemeFolderName != "" {
		themeFolderName = config.Server.ThemeFolderName
	}

	return filepath.Join(config.themeFolderBase, themeFolderName)
}

func (config *Config) load() (*Config, error) {

	path := config.Filepath()

	// check if file can be accessed
	fileInfo, err := os.Open(path)
	if err != nil {
		return config, fmt.Errorf("Cannot read config file %q. Error: %s", path, err)
	}

	defer fileInfo.Close()

	// deserialize config
	serializer := NewJSONSerializer()
	loadedConfig, err := serializer.DeserializeConfig(fileInfo)
	if err != nil {
		return config, fmt.Errorf("Could not deserialize the configuration file %q. Error: %s", path, err)
	}

	// apply values
	config.Server = loadedConfig.Server
	config.Web = loadedConfig.Web

	return config, nil
}

func (config *Config) Save() (*Config, error) {

	path := config.Filepath()

	// make sure the directory exists
	if created, err := util.CreateFile(path); !created {
		return config, fmt.Errorf("Could not create path %q.\nError: %s\n", path, err)
	}

	// open the file for writing
	file, err := os.OpenFile(path, os.O_WRONLY, 0600)
	if err != nil {
		return config, fmt.Errorf("Error while opening file %q for writing.", path)
	}

	writer := bufio.NewWriter(file)

	defer func() {
		writer.Flush()
		file.Close()
	}()

	// serialize the config
	serializer := NewJSONSerializer()
	if serializationError := serializer.SerializeConfig(writer, config); serializationError != nil {
		return config, fmt.Errorf("Error while saving configuration %#v to file %q. Error: %v", config, path, serializationError)
	}

	return config, nil
}

func (config *Config) apply(newConfig *Config) (*Config, error) {
	if newConfig == nil {
		return config, fmt.Errorf("Cannot apply nil.")
	}

	config.Server = newConfig.Server
	config.Web = newConfig.Web

	return config, nil
}

func local(baseFolder string) *Config {
	metaDataFolder := filepath.Join(baseFolder, MetaDataFolderName)
	templatesFolder := filepath.Join(metaDataFolder, TemplatesFolderName)

	return &Config{
		baseFolder:      baseFolder,
		metaDataFolder:  metaDataFolder,
		themeFolderBase: metaDataFolder,
		templatesFolder: templatesFolder,
	}
}

func global(baseFolder string) *Config {
	metaDataFolder := filepath.Join(baseFolder, MetaDataFolderName)
	templatesFolder := filepath.Join(metaDataFolder, TemplatesFolderName)

	return &Config{
		baseFolder:      baseFolder,
		metaDataFolder:  metaDataFolder,
		themeFolderBase: metaDataFolder,
		templatesFolder: templatesFolder,
	}
}

func defaultConfig(baseFolder string) *Config {

	var config *Config

	if isHomeDir(baseFolder) {
		config = global(baseFolder) // global config
	} else {
		config = local(baseFolder) // local config
	}

	// set the default values
	config.Server.ThemeFolderName = ThemeFolderName
	config.Server.Http.Port = 8080
	config.Web.DefaultLanguage = "en"

	return config
}
