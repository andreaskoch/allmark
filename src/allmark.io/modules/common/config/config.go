// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"bufio"
	"fmt"
	"math"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"allmark.io/modules/common/certificates"
	"allmark.io/modules/common/logger/loglevel"
	"allmark.io/modules/common/util/fsutil"
)

const (
	MetaDataFolderName     = ".allmark"
	FilesDirectoryName     = "files"
	ConfigurationFileName  = "config"
	ThemeFolderName        = "theme"
	TemplatesFolderName    = "templates"
	ThumbnailIndexFileName = "thumbnail.index"
	ThumbnailsFolderName   = "thumbnails"
	SSLCertsFolderName     = "certs"

	// Global Defaults
	DefaultHostName                  = "127.0.0.1"
	DefaultHttpPort                  = 0
	DefaultHttpPortEnabled           = true
	DefaultHttpsPort                 = 0
	DefaultHttpsPortEnabled          = true
	DefaultHttpsCertName             = "cert.pem"
	DefaultHttpsKeyName              = "cert.key"
	DefaultForceHttps                = false
	DefaultLanguage                  = "en"
	DefaultLogLevel                  = loglevel.Info
	DefaultIndexingEnabled           = false
	DefaultIndexingIntervalInSeconds = 60
	DefaultLiveReloadEnabled         = false
	DefaultRichTextConversionEnabled = true
)

var homeDirectory func() string

var freeHttpPort int
var freeHttpsPort int

// A flag indicating whether the RTF conversion tool is available
var rtfConversionToolIsAvailable bool

func init() {

	usr, err := user.Current()
	if err != nil {
		panic(fmt.Sprintf("Cannot determine the current users home direcotry. Error: %s", err))
	}

	homeDirectory = func() string {
		return filepath.Clean(usr.HomeDir)
	}

	// check if pandoc is available in the path
	command := exec.Command(DefaultConversionToolPath, "--help")
	if err := command.Run(); err == nil {
		rtfConversionToolIsAvailable = true
	}

	// locate free ports for http and https
	freeHttpPort = getFreePort()
	freeHttpsPort = getFreePort()
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

	return &Config{
		baseFolder:      baseFolder,
		metaDataFolder:  metaDataFolder,
		themeFolderBase: metaDataFolder,
		templatesFolder: templatesFolder,

		Conversion: Conversion{
			Thumbnails: ThumbnailConversion{
				Enabled:       false,
				IndexFileName: ThumbnailIndexFileName,
				FolderName:    ThumbnailsFolderName,
			},
		},
	}
}

func Default(baseFolder string) *Config {

	// create a new config
	config := New(baseFolder)

	// apply default values
	config.Server.ThemeFolderName = ThemeFolderName
	config.Server.Hostname = DefaultHostName
	config.Server.Http.PortNumber = DefaultHttpPort
	config.Server.Http.Enabled = DefaultHttpPortEnabled
	config.Server.Https.PortNumber = DefaultHttpsPort
	config.Server.Https.Enabled = DefaultHttpsPortEnabled
	config.Server.Https.Force = DefaultForceHttps

	config.Web.DefaultLanguage = DefaultLanguage

	// Publisher Information
	config.Web.Publisher = UserInformation{}

	// Authors
	config.Web.Authors = map[string]UserInformation{
		"Unknown": UserInformation{},
	}

	// Thumbnail conversion
	config.Conversion.Thumbnails.IndexFileName = ThumbnailIndexFileName
	config.Conversion.Thumbnails.FolderName = ThumbnailsFolderName

	// Rtf Conversion
	config.Conversion.Rtf.Enabled = DefaultRichTextConversionEnabled

	// Logging
	config.LogLevel = DefaultLogLevel.String()

	// Indexing
	config.Indexing.Enabled = DefaultIndexingEnabled
	config.Indexing.IntervalInSeconds = DefaultIndexingIntervalInSeconds

	// Live-Reload
	config.LiveReload.Enabled = DefaultLiveReloadEnabled

	return config
}

type Port struct {
	PortNumber int
	Enabled    bool
}

func (port *Port) GetPortNumber() int {
	portNumber := port.PortNumber
	if portNumber < 0 || portNumber > math.MaxUint16 {
		panic(fmt.Sprintf("%q is an invalid value for a port. Ports can only be in the range of %v to %v,", portNumber, 0, math.MaxUint16))
	}

	if portNumber == 0 {

		return freeHttpPort

	}

	return portNumber
}

type SecurePort struct {
	Port

	CertFileName string
	KeyFileName  string

	Force bool
}

func (securePort *SecurePort) GetPortNumber() int {
	portNumber := securePort.PortNumber
	if portNumber < 0 || portNumber > math.MaxUint16 {
		panic(fmt.Sprintf("%q is an invalid value for a port. Ports can only be in the range of %v to %v,", portNumber, 0, math.MaxUint16))
	}

	if portNumber == 0 {

		return freeHttpsPort

	}

	return portNumber
}

func (securePort *SecurePort) ForceHttps() bool {
	if securePort.Enabled == false {
		return false
	}

	return securePort.Force
}

type Web struct {
	DefaultLanguage string
	DefaultAuthor   string
	Publisher       UserInformation
	Authors         map[string]UserInformation
}

type UserInformation struct {
	Name  string
	Email string
	Url   string

	GooglePlusHandle string
	TwitterHandle    string
	FacebookHandle   string
}

type Server struct {
	ThemeFolderName string
	Hostname        string
	Http            Port
	Https           SecurePort
}

type Indexing struct {
	Enabled           bool
	IntervalInSeconds int
}

type LiveReload struct {
	Enabled bool
}

type Conversion struct {
	Rtf        RtfConversion
	Thumbnails ThumbnailConversion
}

type RtfConversion struct {
	Enabled bool
}

func (rtf RtfConversion) Tool() string {
	return DefaultConversionToolPath
}

func (rtf RtfConversion) IsEnabled() bool {
	return rtf.Enabled && rtfConversionToolIsAvailable
}

type ThumbnailConversion struct {
	Enabled       bool
	IndexFileName string
	FolderName    string
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
	LiveReload LiveReload
	Analytics  Analytics

	baseFolder      string
	metaDataFolder  string
	themeFolderBase string
	templatesFolder string
}

func (config *Config) CertificateFilePaths() (certificateFilePath, keyFilePath string) {

	// Determine the hostname
	hostname := config.Server.Hostname
	if hostname == "" {
		hostname = DefaultHostName
	}

	// Determine  the cert name
	certificateFileName := config.Server.Https.CertFileName
	if certificateFileName == "" {
		certificateFileName = DefaultHttpsCertName
	}

	// Determine the key name
	keyFileName := config.Server.Https.KeyFileName
	if keyFileName == "" {
		keyFileName = DefaultHttpsKeyName
	}

	// Determine the base directory for the certificates
	certificateBaseDirectory := filepath.Join(config.MetaDataFolder(), SSLCertsFolderName)

	// Default cert and key path
	certificateFilePath = filepath.Join(certificateBaseDirectory, certificateFileName)
	keyFilePath = filepath.Join(certificateBaseDirectory, keyFileName)

	// check if the specified file exists
	if fsutil.FileExists(certificateFilePath) && fsutil.FileExists(keyFilePath) {

		// return the existing paths
		return certificateFilePath, keyFilePath
	}

	// determine the target location for the dummy cert
	if fsutil.DirectoryExists(config.MetaDataFolder()) {

		// meta data folder
		if created := fsutil.CreateDirectory(certificateBaseDirectory); !created {
			panic(fmt.Sprintf("Could not create directory %q", certificateBaseDirectory))
		}

		// the file path can stay the same

	} else {

		// create a temporary directory
		tempDirectory := fsutil.GetTempDirectory()

		certificateFilePath = filepath.Join(tempDirectory, certificateFileName)
		keyFilePath = filepath.Join(tempDirectory, keyFileName)

	}

	// create a dummy cert and key
	if err := certificates.GenerateDummyCert(certificateFilePath, keyFilePath, hostname); err != nil {
		panic(fmt.Sprintf("Could not create dummy certificates %q and %q for hostname %q. Error: %s", certificateFilePath, keyFilePath, hostname, err.Error))
	}

	return certificateFilePath, keyFilePath
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

func (config *Config) ThumbnailIndexFilePath() string {
	filename := ThumbnailIndexFileName
	if config.Conversion.Thumbnails.IndexFileName != "" {
		filename = config.Conversion.Thumbnails.IndexFileName
	}

	return filepath.Join(config.MetaDataFolder(), filename)
}

func (config *Config) ThumbnailFolder() string {
	folderName := ThumbnailsFolderName
	if config.Conversion.Thumbnails.FolderName != "" {
		folderName = config.Conversion.Thumbnails.FolderName
	}

	return filepath.Join(config.MetaDataFolder(), folderName)
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

// Ask the kernel for a free open port that is ready to use
func getFreePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}
