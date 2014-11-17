// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"bufio"
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger/loglevel"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"os"
	"os/user"
	"path/filepath"
)

const (
	MetaDataFolderName    = ".allmark"
	FilesDirectoryName    = "files"
	ConfigurationFileName = "config"
	ThemeFolderName       = "theme"
	TemplatesFolderName   = "templates"
	ThumbnailIndexFile    = "thumbnail.index"
	ThumbnailsFolderName  = "thumbnails"

	// Global Defaults
	DefaultHostName                 = "localhost"
	DefaultPort                     = 8080
	DefaultLanguage                 = "en"
	DefaultLogLevel                 = loglevel.Info
	DefaultReindexIntervalInSeconds = 60
)

var homeDirectory func() string

func init() {

	usr, err := user.Current()
	if err != nil {
		panic(fmt.Sprintf("Cannot determine the current users home direcotry. Error: %s", err))
	}

	homeDirectory = func() string {
		return filepath.Clean(usr.HomeDir)
	}
}

func isHomeDir(directory string) bool {
	return filepath.Clean(directory) == homeDirectory()
}

func Get(baseFolder string) *Config {

	// local
	if config, err := New(baseFolder).Load(); err == nil {
		return config
	}

	// global
	if !isHomeDir(baseFolder) {

		homeDirectory := homeDirectory()
		if config, err := New(homeDirectory).Load(); err == nil {
			return config
		}

	}

	// default
	return Default(baseFolder)
}

func New(baseFolder string) *Config {
	metaDataFolder := filepath.Join(baseFolder, MetaDataFolderName)
	templatesFolder := filepath.Join(metaDataFolder, TemplatesFolderName)

	thumbnailIndexFile := filepath.Join(metaDataFolder, ThumbnailIndexFile)
	thumbnailsFolder := filepath.Join(metaDataFolder, ThumbnailsFolderName)

	return &Config{
		baseFolder:      baseFolder,
		metaDataFolder:  metaDataFolder,
		themeFolderBase: metaDataFolder,
		templatesFolder: templatesFolder,

		thumbnailIndexFile: thumbnailIndexFile,
		thumbnailsFolder:   thumbnailsFolder,
	}
}

func Default(baseFolder string) *Config {

	// create a new config
	config := New(baseFolder)

	// apply default values
	config.Server.ThemeFolderName = ThemeFolderName
	config.Server.Http.Hostname = DefaultHostName
	config.Server.Http.Port = DefaultPort
	config.Web.DefaultLanguage = DefaultLanguage
	config.Conversion.Tool = DefaultConversionToolPath
	config.LogLevel = DefaultLogLevel.String()
	config.Indexing.IntervalInSeconds = DefaultReindexIntervalInSeconds

	return config
}

type Http struct {
	Hostname string
	Port     int
}

type Web struct {
	DefaultLanguage string
}

type Server struct {
	ThemeFolderName string
	Http            Http
}

type Indexing struct {
	IntervalInSeconds int
}

type Conversion struct {
	Tool string
}

type Analytics struct {
	Enabled         bool
	GoogleAnalytics GoogleAnalytics
}

type GoogleAnalytics struct {
	Enabled    bool
	TrackingId string
}

type Config struct {
	Server     Server
	Web        Web
	Conversion Conversion
	LogLevel   string
	Indexing   Indexing
	Analytics  Analytics

	baseFolder      string
	metaDataFolder  string
	themeFolderBase string
	templatesFolder string

	thumbnailIndexFile string
	thumbnailsFolder   string
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

func (config *Config) ThumbnailIndexFile() string {
	return config.thumbnailIndexFile
}

func (config *Config) ThumbnailsFolder() string {
	return config.thumbnailsFolder
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

func (config *Config) Load() (*Config, error) {

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
	config.Conversion = loadedConfig.Conversion
	config.LogLevel = loadedConfig.LogLevel
	config.Indexing = loadedConfig.Indexing
	config.Analytics = loadedConfig.Analytics

	return config, nil
}

func (config *Config) Save() (*Config, error) {

	path := config.Filepath()

	// make sure the directory exists
	if created, err := fsutil.CreateFile(path); !created {
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
	config.Conversion = newConfig.Conversion
	config.LogLevel = newConfig.LogLevel
	config.Indexing = newConfig.Indexing
	config.Analytics = newConfig.Analytics

	return config, nil
}
