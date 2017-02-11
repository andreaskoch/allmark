# Configuration

You can configure and customize how all allmark serves your repositories by creating a custom repository configuration.

Use the `init` action to save the default configuration to the current or given folder:

```bash
cd markdown-repository
allmark init
```

or

```bash
allmark init <directory path>
```

This will create a folder with the name `.allmark` in the current or the specified directory:

- `config`: contains the **JSON configuration** for allmark
- `templates`: contains all **templates** used for rendering the web pages
- `theme`: contains all **assets** used by the templates
- `certs`: contains a generated and self-signed SSL-certificate that can be used for serving HTTPS
- `users.htpasswd`: the user file for **[basic-authentication](http://httpd.apache.org/docs/2.2/programs/htpasswd.html)** (default: `<empty>`)

```
<your-markdown-repository>
└── .allmark
    ├── certs
    │   ├── cert.key
    │   └── cert.pem
    ├── config
    ├── templates
    │   ├── converter.gohtml
    │   ├── document.gohtml
    │   ├── error.gohtml
    │   ├── master.gohtml
    │   ├── opensearchdescription.gohtml
    │   ├── presentation.gohtml
    │   ├── repository.gohtml
    │   ├── robotstxt.gohtml
    │   ├── rssfeed.gohtml
    │   ├── rssfeedcontent.gohtml
    │   ├── search.gohtml
    │   ├── searchcontent.gohtml
    │   ├── sitemap.gohtml
    │   ├── sitemapcontent.gohtml
    │   ├── tagmap.gohtml
    │   ├── tagmapcontent.gohtml
    │   ├── xmlsitemap.gohtml
    │   └── xmlsitemapcontent.gohtml
    ├── theme
    │   ├── autoupdate.js
    │   ├── codehighlighting
    │   │   ├── highlight.css
    │   │   └── highlight.js
    │   ├── deck.css
    │   ├── deck.js
    │   ├── favicon.ico
    │   ├── jquery.js
    │   ├── jquery.lazyload.js
    │   ├── jquery.lazyload.srcset.js
    │   ├── jquery.lazyload.video.js
    │   ├── jquery.tmpl.js
    │   ├── latest.js
    │   ├── modernizr.js
    │   ├── presentation.js
    │   ├── print.css
    │   ├── screen.css
    │   ├── search.js
    │   ├── site.js
    │   ├── tree-last-node.png
    │   ├── tree-node.png
    │   ├── tree-vertical-line.png
    │   └── typeahead.js
    └── users.htpasswd
```

If you init a configuration in your **home-directory**, this configuration will be used as **default for all your repositories** as long as you don't have one in your respective directory:

```bash
allmark init ~/
```

or

```bash
cd ~
allmark init
```

But configurations in your repositories will take precedence over your default configuration in your home-directory.

The **configuration file** has a **JSON format** and is located in `.allmark/config`:

- `Server`
	- `ThemeFolderName`: The name of the folder that contains all theme assets (js, css, ...) (default: `"theme"`)
	- `DomainName`: The default host-/domain name that shall be used (e.g. `"localhost"`, `"www.example.com"`)
	- `HTTP`
		- `Enabled`: If set to `true` http is enabled. If set to `false` http is disabled.
		- `Bindings`: An array of 0..n TCP bindings that will be used to serve HTTP
			- `Network`: `"tcp4"` for IPv4 or `"tcp6"` for IPv6
			- `IP`: An IPv4 address (e.g. `"0.0.0.0"`, `"127.0.0.1"`) or an IPv6 address (e.g. `"::"`, `"::1"`)
			- `Zone`: The [IPv6 zone index](https://en.wikipedia.org/wiki/IPv6_address#Link-local_addresses_and_zone_indices) (e.g. `""`, `"eth0"`, `"eth1"`; (default: `""`)
			- `Port`: 0-65535 (0 means that a random port will be allocated)
	- `HTTPS`
		- `Enabled`: If set to `true` HTTPS is enabled. If set to `false` HTTPS is disabled.
		- `CertFileName`: The filename of the SSL certificate in the `.allmark/certs`-folder (e.g. `"cert.pem"`, `"cert.pem"`)
		- `KeyFileName`: The filename of the SSL certificate key file in the `.allmark/certs`-folder (e.g. `"cert.key"`)
		- `Force`: If set to `true` and if http and HTTPS are enabled all http requests will be redirected to http. If set to `false` you can use HTTPS alongside http.
		- `Bindings`: An array of 0..n TCP bindings that will be used to serve HTTPS
			- same format (Network, IP, Zone, Port) as for HTTP
	- `Authentication`
		- `Enabled`: If set to `true` basic-authentication will be enabled. If set to `false` basic-authentication will be disabled. **Note**: Even if set to `true`, basic authentication will only be enabled if HTTPS is forced.
		- `UserStoreFileName`: The filename of the [htpasswd-file](http://httpd.apache.org/docs/2.2/programs/htpasswd.html) that contains all authorized usernames, realms and passwords/hashes (default: `"users.htpasswd"`).
- `Web`
	- `DefaultLanguage`: An [ISO 639-1](http://en.wikipedia.org/wiki/List_of_ISO_639-1_codes) two-letter language code (e.g. `"en"` → english, `"de"` → german, `"fr"` → french) that is used as the default value for the `<html lang="">` attribute (default: `"en"`).
	- `DefaultAuthor`: The name of the default author (e.g. "John Doe") for all documents in your repository that don't have a `author: Your Name` line in the meta-data section.
	- `Publisher`: Information about the repository-publisher / the owner of an repository.
		- `Name`: The publisher name or organization (e.g. `"Example Org"`)
		- `Email`: The publisher email address (e.g. `"webmaster@example.com"`)
		- `URL`: The URL of the publisher (e.g. `"http://example.com/about"`)
		- `GooglePlusHandle`: The Google+ username/handle of the publisher (e.g. `"exampleorg"`)
		- `TwitterHandle`: The Twitter username/handle of the publisher (e.g. `"exampleorg"`)
		- `FacebookHandle`: The Facebook username/handle of the publisher (e.g. `"exampleorg"`)
	- `Authors`: Detail information about each author of your repository (by name)
		- `"John Doe"`
			- `Name`: `"John Doe"`
			- `Email`: `"johndoe@example.com"`
			- `URL`: `"http://example.com/about/johndoe"`
			- `GooglePlusHandle`: The Google+ username/handle of the author (e.g. `"johndoe"`)
			- `TwitterHandle`: The Twitter username/handle of the author (e.g. `"johndoe"`)
			- `FacebookHandle`: The Facebook username/handle of the author (e.g. `"johndoe"`)
		- `"Jane Doe"`
			- `"Name"`
			- ...
		- ...
- `Conversion`
	- `RTF`: Rich-text Conversion
		- `Enabled`: If set to `true` rich-text conversion is enabled. allmark uses [pandoc](http://pandoc.org/) for the rich-text conversion. If the [pandoc binary](https://github.com/jgm/pandoc/releases/latest) is not found in your PATH, rich-text conversion will not be available.
	- `Thumbnails`: Image-Thumbnail creation.
		- `Enabled`: If set to `true` allmark will create smaller versions (Small: 320x240, Medium: 640x480, Large: 1024x768) for all images in your repository and use the respective version depending on the screen size of your clients (default: `false`).
	- `IndexFileName`: The name of the file where allmark stores an index of all thumbnails it has created (default: `"thumbnail.index"`).
	- `FolderName`: The name of the folder were allmark stores the thumbnails (default: `"thumbnails"`).
- `LogLevel`: Possible options are: `"off"`, `"debug"`, `"info"`, `"statistics"`, `"warn"`, `"error"`, `"fatal"` (default: `"info"`).
- `Indexing`
	- `IntervalInSeconds`: The indexing interval in seconds (default: 60). allmark will reindex the repository every x seconds.
- `Analytics`
	- `Enabled`: If set to `true` analytics is enabled (default: `false`).
	- `GoogleAnalytics`
		- `Enabled`: If set to `true` Google Analytics is enabled (default: `false`).
		- `TrackingID`: Your Google Analytics tracking id (e.g `"UA-000000-01"`).


```json
{
	"Server": {
		"ThemeFolderName": "theme",
		"DomainName": "localhost",
		"HTTP": {
			"Enabled": true,
			"Bindings": [
				{
					"Network": "tcp4",
					"IP": "0.0.0.0",
					"Zone": "",
					"Port": 80
				},
				{
					"Network": "tcp6",
					"IP": "::",
					"Zone": "",
					"Port": 80
				}
			]
		},
		"HTTPS": {
			"Enabled": true,
			"Bindings": [
				{
					"Network": "tcp4",
					"IP": "0.0.0.0",
					"Zone": "",
					"Port": 443
				},
				{
					"Network": "tcp6",
					"IP": "::",
					"Zone": "",
					"Port": 443
				}
			],
			"CertFileName": "cert.pem",
			"KeyFileName": "cert.key",
			"Force": false
		},
		"Authentication": {
			"Enabled": false,
			"UserStoreFileName": "users.htpasswd"
		}
	},
	"Web": {
		"DefaultLanguage": "en",
		"DefaultAuthor": "",
		"Publisher": {
			"Name": "",
			"Email": "",
			"URL": "",
			"GooglePlusHandle": "",
			"TwitterHandle": "",
			"FacebookHandle": ""
		},
		"Authors": {
			"Unknown": {
				"Name": "",
				"Email": "",
				"URL": "",
				"GooglePlusHandle": "",
				"TwitterHandle": "",
				"FacebookHandle": ""
			}
		}
	},
	"Conversion": {
		"RTF": {
			"Enabled": true
		},
		"Thumbnails": {
			"Enabled": false,
			"IndexFileName": "thumbnail.index",
			"FolderName": "thumbnails"
		}
	},
	"LogLevel": "Info",
	"Indexing": {
		"IntervalInSeconds": 60
	},
	"Analytics": {
		"Enabled": false,
		"GoogleAnalytics": {
			"Enabled": false,
			"TrackingID": ""
		}
	}
}
```

---

created at: 2015-08-03
modified at: 2015-08-03
author: Andreas Koch
tags: Configuration, Documentation
alias: configuration
