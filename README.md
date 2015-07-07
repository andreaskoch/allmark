# allmark - the markdown server

allmark is a standalone markdown web server for Linux, Mac OS and Windows written in go.

![allmark logo (128x128px)](files/design/logo/PNG8/allmark-logo-128x128.png)

## Usage

Serve a specific directory:

```bash
allmark serve <directory path>
```

Serve the current directory:

```bash
cd markdown-repository
allmark serve
```

Serve the current directory with **live-reload** enabled:

```bash
allmark serve -livereload
```

Force a full **reindex** every 60* seconds:

```bash
allmark serve -reindex
```

`*` The default interval is 60 seconds. You can change the interval in the repository config.

Force HTTPs (redirect all http requests to https):

```bash
allmark serve -secure
```

Save the default configuration to the `.allmark` folder so you can customize it:

```bash
allmark init
```

You can point **allmark** at any folder structure that contains **markdown documents** and files referenced by these documents (e.g. this repository folder) and allmark will start a **web-server** and serve the folder contents as HTML via HTTP(s) on a random free port.

**Folder Structure Conventions**

The standard folder structure for a **markdown-repository item** could look something like this:

```
├── files
│   ├── image.png
│   └── more-files
│       ├── file1.txt
│       ├── file2.txt
│       └── file3.txt
└── some-file.md
```

1. one markdown file per folder (with the extension .md, .markdown or .mdown)
2. a `files` folder which contains all files referenced by the markdown document
3. an arbitrary number of child directories that can contain more markdown-repository items

**Nesting / Hierarchie**

You can nest repository items arbitrarily. Example:

```
├── child-item-1
│   └── item1.md
├── child-item-2
│   └── item2.md
├── child-item-3
│   └── item3.md
├── files
│   ├── image.png
│   └── more-files
│       ├── file1.txt
│       ├── file2.txt
│       └── file3.txt
└── some-file.md
```

**Folders without Markdown Files**

- If you have folders in your repository that don't contains markdown files allmark will display and index of all files in that directory (→ **file-collection item**)
- file-collection items cannot have other childs

**Markdown Document Structure**

allmark makes certain assumptions about the structure of your documents. They should have

1. Title
2. Description Text
3. Document Body

A **typical document** expected by allmark could look like this:

	# Document Title / Headline

	A short description of the document ... Usually one sentence.

	The Content of your document

	![Some Image](files/image.jpg)

	- A List 1
	- A List 2
	- A List 3

	**Some garbage text**: In pharetra ullamcorper egestas.
	Nam vel sodales velit. Nulla elementum dapibus sem nec scelerisque.
	In hac habitasse platea dictumst. Nulla vestibulum lacinia tincidunt.

## Download / Installation

You can download the **latest binaries** of allmark for your operating system from [allmark.io/bin](https://allmark.io/bin)

**Linux**

```bash
sudo su
curl -s --insecure https://allmark.io/bin/allmark > /usr/local/bin/allmark
chmod +x /usr/local/bin/allmark
```

**Mac OS**

```bash
sudo curl "https://allmark.io/bin/darwin_amd64/allmark" -o "/usr/local/bin/allmark"
sudo chmod +x /usr/local/bin/allmark
```

**Windows**

```powershell
Invoke-WebRequest https://allmark.io/bin/windows_amd64/allmark.exe -OutFile allmark.exe
```

All binaries at [allmark.io](https://allmark.io) are up-to-date builds of the **master**-branch.

If you want to download and install binaries from the **develop**-branch you can go to [develop.allmark.io/bin](https://develop.allmark.io).

## Configuration

You can configure and customize how all allmark serves your repositories by creating a custom repository configuration.

You can use the `init` action to save the default configuration to the current or given folder:

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
- `certs`: contains a generated and self-signed SSL-certificate that can be used for serving HTTPs
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
    │   ├── pdf.js
    │   ├── pdfpreview.js
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
	- `HTTPs`
		- `Enabled`: If set to `true` https is enabled. If set to `false` https is disabled.
		- `CertFileName`: The filename of the SSL certificate in the `.allmark/certs`-folder (e.g. `"cert.pem"`, `"cert.pem"`)
		- `KeyFileName`: The filename of the SSL certificate key file in the `.allmark/certs`-folder (e.g. `"cert.key"`)
		- `Force`: If set to `true` and if http and https are enabled all http requests will be redirected to http. If set to `false` you can use https alongside http.
		- `Bindings`: An array of 0..n TCP bindings that will be used to serve HTTPs
			- same format (Network, IP, Zone, Port) as for HTTP
	- `Authentication`
		- `Enabled`: If set to `true` basic-authentication will be enabled. If set to `false` basic-authentication will be disabled. **Note**: Even if set to `true`, basic authentication will only be enabled if HTTPs is forced.
		- `UserStoreFileName`: The filename of the [htpasswd-file](http://httpd.apache.org/docs/2.2/programs/htpasswd.html) that contains all authorized usernames, realms and passwords/hashes (default: `"users.htpasswd"`).
- `Web`
	- `DefaultLanguage`: An [ISO 639-1](http://en.wikipedia.org/wiki/List_of_ISO_639-1_codes) two-letter language code (e.g. `"en"` → english, `"de"` → german, `"fr"` → french) that is used as the default value for the `<html lang="">` attribute (default: `"en"`).
	- `DefaultAuthor`: The name of the default author (e.g. "John Doe") for all documents in your repository that don't have a `author: Your Name` line in the meta-data section.
	- `Publisher`: Information about the repository-publisher / the owner of an repository.
		- `Name`: The publisher name or organization (e.g. `"Example Org"`)
		- `Email`: The publisher email address (e.g. `"webmaster@example.com"`)
		- `Url`: The URL of the publisher (e.g. `"http://example.com/about"`)
		- `GooglePlusHandle`: The Google+ username/handle of the publisher (e.g. `"exampleorg"`)
		- `TwitterHandle`: The Twitter username/handle of the publisher (e.g. `"exampleorg"`)
		- `FacebookHandle`: The Facebook username/handle of the publisher (e.g. `"exampleorg"`)
	- `Authors`: Detail information about each author of your repository (by name)
		- `"John Doe"`
			- `Name`: `"John Doe"`
			- `Email`: `"johndoe@example.com"`
			- `Url`: `"http://example.com/about/johndoe"`
			- `GooglePlusHandle`: The Google+ username/handle of the author (e.g. `"johndoe"`)
			- `TwitterHandle`: The Twitter username/handle of the author (e.g. `"johndoe"`)
			- `FacebookHandle`: The Facebook username/handle of the author (e.g. `"johndoe"`)
		- `"Jane Doe"`
			- `"Name"`
			- ...
		- ...
- `Conversion`
	- `Rtf`: Rich-text Conversion
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
		- `TrackingId`: Your Google Analytics tracking id (e.g `"UA-000000-01"`).


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
		"HTTPs": {
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
			"Url": "",
			"GooglePlusHandle": "",
			"TwitterHandle": "",
			"FacebookHandle": ""
		},
		"Authors": {
			"Unknown": {
				"Name": "",
				"Email": "",
				"Url": "",
				"GooglePlusHandle": "",
				"TwitterHandle": "",
				"FacebookHandle": ""
			}
		}
	},
	"Conversion": {
		"Rtf": {
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
			"TrackingId": ""
		}
	}
}
```

## Features

This is an unordered list of the most prominent features of allmark:

1. Renders [GitHub Flavored MarkDown](https://help.github.com/articles/github-flavored-markdown/)
2. Full text search (+ Autocomplete)
3. Live-Reload / Live-Editing (via WebSockets)
4. Document Tagging
5. Tag Cloud
6. Documents By Tag
7. HTML Sitemap
8. XML Sitemap
9. robots.txt
10. RSS Feed
11. Print Preview
12. JSON Representation of Documents
13. Hierarchical Document Trees
14. Repository Navigation
	- Top-Level Navigation
	- Bread-Crumb Navigation
	- Previous and Next Items
	- Child-Documents
15. Image Thumbnails
16. Markdown Extensions
	- Image Galleries
	- File Preview
	- Displaying Folder Contents
	- Video Player Integration
	- Audio Player Integration
	- PDF Document Preview
	- Repository cross-links by alias
17. Different Item Types (Repository, Document, Presentation)
18. Document Meta Data
	- Author
	- Tags
	- Document Alias
	- Creation Date
	- Last Modified Date
	- Language
	- Geo Location
19. Default Theme
	- Responsive Design
	- Lazy Loading for images and videos
	- Syntax Highlighting
20. Presentation Mode
21. Rich Text Conversion (Download documents as .rtf files)
22. Image Thumbnail Generation
23. HTTPs Support
	- Reference custom SSL certificates via `.allmark/config` from the `.allmark/certs` folder
	- Generates self-signed SSL certificates on-the-fly if no certificate is configured
24. Basic-Authentication
	- For an additional level of security allmark will only allow basic-authentication over SSL.
	- You can add users to the `.allmark/users.htpasswd` file using the tool [htpasswd](http://httpd.apache.org/docs/2.2/programs/htpasswd.html)
25. Parallel hosting of HTTP/HTTPs over IPv4 and/or IPv6

I will try to create videos showing you the different features when there is time.

## Demo / Showcase

If you want to see **allmark in action** you can visit my blog [AndyK Docs](https://andykdocs.de/) at [https://andykdocs.de](https://andykdocs.de):

![Animation: Demo of allmark hosting andykdocs.de ](files/demo/allmark-demo-andykdocs.gif)

## Build Status

[![Build Status](https://travis-ci.org/andreaskoch/allmark.png)](https://travis-ci.org/andreaskoch/allmark)

There is also an automated docker build at [registry.hub.docker.com/u/andreaskoch/allmark/](https://registry.hub.docker.com/u/andreaskoch/allmark/) which builds the develop and master branch every time a commit is pushed.

## Build

If you have [go](https://golang.org/dl/) (≥ 1.3) installed you can build allmark yourself in two steps:

1. Clone the project from github
2. Run the `make.go` file with the `-install` flag

```bash
git clone git@github.com:andreaskoch/allmark.git
cd allmark
go run make.go -install
```

Afterwards you will find the `allmark` binary in the bin-folder of the project. To test your installation you can start by serving the allmark-project directory:

```bash
cd allmark
bin/allmark serve
```

After a second or so a browser window should pop up.

![Screenshot: Testing the allmark server on the allmark-project directory](files/installation/screenshot-allmark-test-run-on-project-folder.png)

## Dependencies

allmark relies on a number of third-party libraries:

- [github.com/bradleypeabody/fulltext](src/github.com/bradleypeabody/fulltext)
- [github.com/gorilla/context](src/github.com/gorilla/context)
- [github.com/gorilla/mux](src/github.com/gorilla/mux)
- [github.com/gorilla/handlers](src/github.com/gorilla/handlers)
- [github.com/jbarham/go-cdb](src/github.com/jbarham/go-cdb)
- [github.com/nfnt/resize](src/github.com/nfnt/resize)
- [github.com/russross/blackfriday](src/github.com/russross/blackfriday)
- [github.com/shurcooL/go/github_flavored_markdown/sanitized_anchor_name](src/github.com/shurcooL/go/github_flavored_markdown/sanitized_anchor_name)
- [github.com/skratchdot/open-golang/open](src/github.com/skratchdot/open-golang/open)
- [golang.org/x/net/websocket](src/golang.org/x/net/websocket)
- [github.com/andreaskoch/go-fswatch](src/github.com/andreaskoch/go-fswatch)
- [github.com/abbot/go-http-auth](src/github.com/abbot/go-http-auth)

These dependencies are not covered by the allmark copyright/license. See the respective projects for their copyright & licensing details.

The packages are mirrored into allmark [src-folder](src) for hermetic build reasons and versioning.

To get a full list of all used third-party libraries you can execute the make tool with the `-list-dependencies` flag:

```bash
go run make.go -list-dependencies
```

To download the latest versions for all dependencies use the `-update-dependencies` flag:

```bash
go run make.go -update-dependencies
```

## Cross-Compilation

If you want to cross-compile allmark for different platforms (darwin, dragonfly, freebsd, linux, nacl, netbsd, openbsd, windows) and architectures (386, amd64, arm) you can do so by using the `-crosscompile` or `-crosscompile-with-docker` flags for the make script.

If you have prepared your go environment for cross-compilation (see: [Dave Cheney - An introduction to cross compilation with Go](http://dave.cheney.net/2012/09/08/an-introduction-to-cross-compilation-with-go)) you can use the `-crosscompile` flag:

```bash
go run make.go -crosscompile
```

This command will cross-compile for all platforms and architectures directly on your system.

If you have not prepared your golang installation for cross-compilation you can use the the `-crosscompile-with-docker` flag instead. This command will launch a [docker container with go 1.4](https://registry.hub.docker.com/u/library/golang/) that is prepared for cross-compilation and will build allmark for all available platforms and architectures inside the docker-container. The output will be available in the `bin` folder of this project:

```
bin/
├── allmark
..
├── darwin_amd64
│   └── allmark
..
├── dragonfly_amd64
│   └── allmark
..
├── freebsd_amd64
│   └── allmark
├── freebsd_arm
│   └── allmark
..
├── linux_arm
│   └── allmark
..
├── nacl_arm
│   └── allmark
..
├── netbsd_amd64
│   └── allmark
..
├── openbsd_amd64
│   └── allmark
..
├── windows_amd64
│   └── allmark.exe
├── ...
├── README.md
└── src
```

## Known Bugs

### Windows

- Filesystem links: Serving folders that are filesystem junctions/links is no longer possible with go 1.4 (it did work with go 1.3)

## Roadmap / To Dos

Here are some of the ideas and todos I would like to add in the future. Contributions are welcome!

### Architecture & Features

- Expose the markdown source
- Web Editor for Markdown Documents
- Data Access
    - Dropbox support
    - SMTP message posting
    - Repository Replication?
    - Amazon S3
- allmark swarm
    - Repository sharding
    - load-balancing
- Static website generation
- User Management / Access Restrictions
    - User management pages
- Support for Folders with multiple Markdown Files
- Support for custom-rewrites

### Theming

- Redesign default theme with Twitter Bootstrap
- Create a theme "loader"
- Infinite Scrolling for latest items
- Improved Image Galleries
