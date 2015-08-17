// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package config provides access to the configuration models required by
// the rest of the modules.
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
	"allmark.io/modules/common/ports"
	"allmark.io/modules/common/util/fsutil"
	"github.com/abbot/go-http-auth"
)

// Global configuration constants.
const (
	MetaDataFolderName     = ".allmark"
	FilesDirectoryName     = "files"
	ConfigurationFileName  = "config"
	ThemeFolderName        = "theme"
	TemplatesFolderName    = "templates"
	ThumbnailIndexFileName = "thumbnail.index"
	ThumbnailsFolderName   = "thumbnails"
	SSLCertsFolderName     = "certs"
)

// Global default values.
const (
	DefaultDomainName                = "localhost"
	DefaultHTTPPortEnabled           = true
	DefaultHTTPSPortEnabled          = false
	DefaultHTTPSCertName             = "cert.pem"
	DefaultHTTPSKeyName              = "cert.key"
	DefaultForceHTTPS                = false
	DefaultLanguage                  = "en"
	DefaultLogLevel                  = loglevel.Error
	DefaultIndexingEnabled           = false
	DefaultIndexingIntervalInSeconds = 60
	DefaultLiveReloadEnabled         = false
	DefaultRichTextConversionEnabled = true
	DefaultAuthenticationEnabled     = false
	DefaultUserStoreFileName         = "users.htpasswd"
)

var homeDirectory func() string

// A flag indicating whether the DOCX conversion tool is available
var docxConversionToolIsAvailable bool

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
		docxConversionToolIsAvailable = true
	}
}

func isHomeDir(directory string) bool {
	return filepath.Clean(directory) == homeDirectory()
}

// Get tries to locate a Config in the specified baseFolder and return it.
// If no configuration was found in the specified folder it will check the users home-directory for a Config.
// If no configuration was found in the supplied baseFolder and no global config in the users home-directory Get will return a default configuration.
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

// New creates a new configuration for the given baseFolder.
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

// Default returns a configuration with default values for the given baseFolder.
func Default(baseFolder string) *Config {

	// create a new config
	config := New(baseFolder)

	// apply default values
	config.Server.ThemeFolderName = ThemeFolderName
	config.Server.DomainName = DefaultDomainName

	// HTTP
	config.Server.HTTP.Enabled = DefaultHTTPPortEnabled
	config.Server.HTTP.Bindings = []*TCPBinding{
		&TCPBinding{
			Network: "tcp4",
			IP:      "0.0.0.0",
			Zone:    "",
			Port:    0,
		},
		&TCPBinding{
			Network: "tcp6",
			IP:      "::",
			Zone:    "",
			Port:    0,
		},
	}

	// HTTPS
	config.Server.HTTPS.Enabled = DefaultHTTPSPortEnabled
	config.Server.HTTPS.CertFileName = DefaultHTTPSCertName
	config.Server.HTTPS.KeyFileName = DefaultHTTPSKeyName
	config.Server.HTTPS.Force = DefaultForceHTTPS
	config.Server.HTTPS.Bindings = []*TCPBinding{
		&TCPBinding{
			Network: "tcp4",
			IP:      "0.0.0.0",
			Zone:    "",
			Port:    0,
		},
		&TCPBinding{
			Network: "tcp6",
			IP:      "::",
			Zone:    "",
			Port:    0,
		},
	}

	// Authentication
	config.Server.Authentication.Enabled = DefaultAuthenticationEnabled
	config.Server.Authentication.UserStoreFileName = DefaultUserStoreFileName

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

	// DOCX Conversion
	config.Conversion.DOCX.Enabled = DefaultRichTextConversionEnabled

	// Logging
	config.LogLevel = DefaultLogLevel.String()

	// Indexing
	config.Indexing.Enabled = DefaultIndexingEnabled
	config.Indexing.IntervalInSeconds = DefaultIndexingIntervalInSeconds

	// Live-Reload
	config.LiveReload.Enabled = DefaultLiveReloadEnabled

	return config
}

// TCPBinding contains all required parameters for a tcp4 or tcp6 address binding.
type TCPBinding struct {
	Network string

	IP   string
	Zone string
	Port int
}

// GetTCPAddress returns a net.TCPAddress object of the current TCP binding.
func (binding *TCPBinding) GetTCPAddress() net.TCPAddr {
	ip := net.ParseIP(binding.IP)
	return net.TCPAddr{
		IP:   ip,
		Port: binding.Port,
		Zone: binding.Zone,
	}
}

// AssignFreePort locates a free port and assigns it the the current binding.
func (binding *TCPBinding) AssignFreePort() {
	if binding.Port > 0 && binding.Port < math.MaxUint16 {
		return
	}

	binding.Port = ports.GetFreePort(binding.Network, binding.GetTCPAddress())
}

// HTTP contains the configuration parameters for HTTP server endpoint.
type HTTP struct {
	Enabled  bool
	Bindings []*TCPBinding
}

// HTTPS contains the configuration parameters for a HTTPS server endpoint.
type HTTPS struct {
	HTTP

	CertFileName string
	KeyFileName  string

	Force bool
}

// HTTPSIsForced indicates whether HTTPS is forced or not.
func (https *HTTPS) HTTPSIsForced() bool {
	if https.Enabled == false {
		return false
	}

	return https.Force
}

// Authentication contains authentication settings.
type Authentication struct {
	// Enabled is flag indicating whether authentication is enabled.
	Enabled bool

	// UserStoreFileName defines the file name for the authentication user-store file (e.g. "users.htpasswd").
	UserStoreFileName string
}

// Web contains all web-site related properties such as the language, authors and publisher information.
type Web struct {
	DefaultLanguage string
	DefaultAuthor   string
	Publisher       UserInformation
	Authors         map[string]UserInformation
}

// UserInformation contains user-related properties such as the Name and Email address.
type UserInformation struct {
	Name  string
	Email string
	URL   string

	GooglePlusHandle string
	TwitterHandle    string
	FacebookHandle   string
}

// Server contains web-server related parameters such as the domain-name, theme-folder and HTTP/HTTPs bindings.
type Server struct {
	ThemeFolderName string
	DomainName      string
	HTTP            HTTP
	HTTPS           HTTPS
	Authentication  Authentication
}

// Indexing defines the reindexing parameters of the repository.
type Indexing struct {
	Enabled           bool
	IntervalInSeconds int
}

// LiveReload defines the live-reload capabilities.
type LiveReload struct {
	Enabled bool
}

// Conversion defines the rich-text and thumbnail conversion paramters.
type Conversion struct {
	DOCX       DOCXConversion
	Thumbnails ThumbnailConversion
}

// DOCXConversion contains rich-text (DOCX) conversion parameters.
type DOCXConversion struct {
	Enabled bool
}

// Tool returns the path of the external rich-text conversion tool (pandoc) used
// to create Rich-text documents from repository items.
func (docx DOCXConversion) Tool() string {
	return DefaultConversionToolPath
}

// IsEnabled returns a flag indicating if rich-text conversion is enabled or not.
// Rich-text conversion can only be enabled if the conversion tool was found in the PATH on startup.
func (docx DOCXConversion) IsEnabled() bool {
	return docx.Enabled && docxConversionToolIsAvailable
}

// ThumbnailConversion defines the image-thumbnail conversion capabilities.
type ThumbnailConversion struct {
	Enabled       bool
	IndexFileName string
	FolderName    string
}

// Analytics defines the web-analytics parameters of the web-server.
type Analytics struct {
	Enabled         bool
	GoogleAnalytics GoogleAnalytics
}

// GoogleAnalytics contains the Google Analytics realted parameters for the web-analytics section.
type GoogleAnalytics struct {
	Enabled    bool
	TrackingID string
}

// Config is the main configuration model for all parts of allmark.
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

// CertificateDirectory returns the path of the SSL certificates directory in the meta-data folder.
func (config *Config) CertificateDirectory() string {
	return filepath.Join(config.MetaDataFolder(), SSLCertsFolderName)
}

// CertificateFilePaths returns the SSL certificate and key file paths.
// If none are configured or the configured ones don't exist it will create new
// ones and return the paths of the newly generates certificate/key pair.
func (config *Config) CertificateFilePaths() (certificateFilePath, keyFilePath string) {

	// Determine the domain name
	domainname := config.Server.DomainName
	if domainname == "" {
		domainname = DefaultDomainName
	}

	// Determine  the cert name
	certificateFileName := config.Server.HTTPS.CertFileName
	if certificateFileName == "" {
		certificateFileName = DefaultHTTPSCertName
	}

	// Determine the key name
	keyFileName := config.Server.HTTPS.KeyFileName
	if keyFileName == "" {
		keyFileName = DefaultHTTPSKeyName
	}

	// Determine the base directory for the certificates
	certificateBaseDirectory := config.CertificateDirectory()

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
	if err := certificates.GenerateDummyCert(certificateFilePath, keyFilePath, domainname); err != nil {
		panic(fmt.Sprintf("Could not create dummy certificates %q and %q for domainname %q. Error: %s", certificateFilePath, keyFilePath, domainname, err.Error()))
	}

	return certificateFilePath, keyFilePath
}

// BaseFolder returns the path of the base folder of the current configuration model.
func (config *Config) BaseFolder() string {
	return config.baseFolder
}

// MetaDataFolder returns the path of the meta-data folder.
func (config *Config) MetaDataFolder() string {
	return config.metaDataFolder
}

// TemplatesFolder returns the path of the templates folder.
func (config *Config) TemplatesFolder() string {
	return config.templatesFolder
}

// Filepath returns the path of the serialized version of the current configuration model.
func (config *Config) Filepath() string {
	return filepath.Join(config.MetaDataFolder(), ConfigurationFileName)
}

// ThemeFolder returns the path of the theme folder.
func (config *Config) ThemeFolder() string {
	themeFolderName := ThemeFolderName
	if config.Server.ThemeFolderName != "" {
		themeFolderName = config.Server.ThemeFolderName
	}

	return filepath.Join(config.themeFolderBase, themeFolderName)
}

// ThumbnailIndexFilePath returns the path of the thumbnail index file.
func (config *Config) ThumbnailIndexFilePath() string {
	filename := ThumbnailIndexFileName
	if config.Conversion.Thumbnails.IndexFileName != "" {
		filename = config.Conversion.Thumbnails.IndexFileName
	}

	return filepath.Join(config.MetaDataFolder(), filename)
}

// ThumbnailFolder returns the path of the thumbnail folder.
func (config *Config) ThumbnailFolder() string {
	folderName := ThumbnailsFolderName
	if config.Conversion.Thumbnails.FolderName != "" {
		folderName = config.Conversion.Thumbnails.FolderName
	}

	return filepath.Join(config.MetaDataFolder(), folderName)
}

// Load reads the configuration-model from disk.
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

// Save persists the current state of the configuration-model to disk.
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

// AuthenticationIsEnabled get a flag indicating if authentication is enabled.
func (config *Config) AuthenticationIsEnabled() bool {

	if !config.Server.Authentication.Enabled {
		return false
	}

	// we will only allow basic authentication over https
	if config.Server.HTTP.Enabled && config.Server.HTTPS.HTTPSIsForced() == false {
		fmt.Println("Basic-Authentication over HTTP is not available. Please disable HTTP or force HTTPS in order to use basic-authentication.")
		os.Exit(1)
	}

	return true
}

// AuthenticationFilePath returns the path of the authentication file.
func (config *Config) AuthenticationFilePath() string {

	if config.Server.Authentication.UserStoreFileName == "" {
		config.Server.Authentication.UserStoreFileName = DefaultUserStoreFileName
	}

	digestFilePath := filepath.Join(config.MetaDataFolder(), config.Server.Authentication.UserStoreFileName)
	return digestFilePath
}

// GetAuthenticationUserStore returns a digest-access authentication secret provider function
// that uses the configured authentication file.
func (config *Config) GetAuthenticationUserStore() auth.SecretProvider {
	// abort if authentication is disabled
	if !config.AuthenticationIsEnabled() {
		return nil
	}

	// panic if authentication is enabled but the auth file does not exist.
	digestFilePath := config.AuthenticationFilePath()
	if !fsutil.FileExists(digestFilePath) {
		panic(fmt.Sprintf("The specified authentication user store %q does not exist.", digestFilePath))
	}

	return auth.HtpasswdFileProvider(digestFilePath)
}
