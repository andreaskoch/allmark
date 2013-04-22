// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"bufio"
	"fmt"
	"github.com/andreaskoch/allmark/themes"
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

func Initialize(baseFolder string) (*Config, error) {
	homeDirectory, homeDirError := getUserHomeDir()
	if homeDirError != nil {
		return nil, fmt.Errorf("Cannot determine the current users home directory location.")
	}

	if filepath.Clean(baseFolder) == filepath.Clean(homeDirectory) {
		return initializeGlobal(baseFolder)
	}

	return initializeLocal(baseFolder)
}

func GetConfig(repositoryPath string) *Config {

	// return the local config
	if exists, localConfig := getLocalConfig(repositoryPath); exists {
		return localConfig
	}

	// return the global config
	if homeDirectory, homeDirError := getUserHomeDir(); homeDirError == nil {
		if exists, globalConfig := getGlobalConfig(homeDirectory); exists {
			return globalConfig
		}
	}

	// return the default config
	return defaultConfig(repositoryPath)
}

func createTheme(baseFolder string) (success bool, err error) {
	if !util.CreateDirectory(baseFolder) {
		return false, fmt.Errorf("Unable to create theme folder %q.", baseFolder)
	}

	themeFile := filepath.Join(baseFolder, "screen.css")
	file, err := os.Create(themeFile)
	if err != nil {
		return false, fmt.Errorf("Unable to create theme file %q.", themeFile)
	}

	defer file.Close()
	file.WriteString(themes.GetTheme())

	return true, nil
}

func initializeLocal(baseFolder string) (*Config, error) {

	// get the existing configuration
	exists, existingConfig := getLocalConfig(baseFolder)
	if !exists {
		existingConfig = defaultConfig(baseFolder)
	}

	// create a new configuration
	config := local(baseFolder)
	config.apply(existingConfig)

	// create config
	if _, err := config.save(); err != nil {
		return nil, fmt.Errorf("Error while creating configuration file %q. Error: ", config.Filepath(), err)
	}

	fmt.Printf("Local configuration created at %q.\n", config.Filepath())

	// create theme
	themeFolder := config.ThemeFolder()
	if success, err := createTheme(themeFolder); !success {
		return nil, fmt.Errorf("%s", err)
	}

	fmt.Printf("Local theme created at %q.\n", themeFolder)

	return config, nil
}

func initializeGlobal(baseFolder string) (*Config, error) {

	// get the existing configuration
	exists, existingConfig := getGlobalConfig(baseFolder)
	if !exists {
		existingConfig = defaultConfig(baseFolder)
	}

	// create a new configuration
	config := global(baseFolder)
	config.apply(existingConfig)

	// create config
	if _, err := config.save(); err != nil {
		return nil, fmt.Errorf("Error while creating configuration file %q. Error: ", config.Filepath(), err)
	}

	fmt.Printf("Global configuration created at %q.\n", config.Filepath())

	// create theme
	themeFolder := config.ThemeFolder()
	if success, err := createTheme(themeFolder); !success {
		return nil, fmt.Errorf("%s", err)
	}

	fmt.Printf("Global theme created at %q.\n", themeFolder)

	return config, nil
}

func getLocalConfig(baseFolder string) (exists bool, config *Config) {
	if config, err := local(baseFolder).load(); err == nil {
		return true, config
	}

	return false, nil
}

func getGlobalConfig(baseFolder string) (exists bool, config *Config) {
	if config, err := global(baseFolder).load(); err == nil {
		return true, config
	}

	return false, nil
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
}

func (config *Config) BaseFolder() string {
	return config.baseFolder
}

func (config *Config) MetaDataFolder() string {
	return config.metaDataFolder
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

func (config *Config) save() (*Config, error) {

	path := config.Filepath()

	// create or overwrite the config file
	if created, err := util.CreateFile(path); !created {
		return config, fmt.Errorf("Could not create configuration file %q. Error: ", path, err)
	}

	// open the file for writing
	file, err := os.OpenFile(path, os.O_WRONLY, 0776)
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

	return &Config{
		baseFolder:      baseFolder,
		metaDataFolder:  metaDataFolder,
		themeFolderBase: baseFolder,
	}
}

func global(baseFolder string) *Config {
	metaDataFolder := filepath.Join(baseFolder, MetaDataFolderName)

	return &Config{
		baseFolder:      baseFolder,
		metaDataFolder:  metaDataFolder,
		themeFolderBase: metaDataFolder,
	}
}

func defaultConfig(baseFolder string) *Config {
	defaultConfig := local(baseFolder)
	defaultConfig.Server.ThemeFolderName = ThemeFolderName
	defaultConfig.Server.Http.Port = 8080
	defaultConfig.Web.DefaultLanguage = "en"

	return defaultConfig
}

// Get the current users home directory path
func getUserHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("Cannot determine the current users home direcotry. Error: %s", err)
	}

	return usr.HomeDir, nil
}
