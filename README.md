# allmark - the markdown server

allmark is a standalone markdown web server for Linux, BSD, Solaris Mac OS and Windows written in go.

![allmark logo (128x128px)](files/design/logo/PNG8/allmark-logo-128x128.png)

You can point **allmark** at any folder structure that contains **markdown documents** and files referenced by these documents (e.g. this repository folder) and allmark will start a **web-server** and serve the folder contents as HTML via HTTP (default: 8080).

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

- If you have folders in your repository that don't contains markdown files allmark will display and index of all files in that directory (=> **file-collection item**)
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
sudo curl -s --insecure https://allmark.io/bin/allmark > /usr/local/bin/allmark
sudo chmod +x /usr/local/bin/allmark
```

**Mac OS**

```bash
sudo curl "https://allmark.io/bin/darwin_amd64/allmark" -o "/usr/local/bin/allmark"
sudo chmod +x /usr/local/bin/allmark
```

**Windows**

```bash
Invoke-WebRequest https://allmark.io/bin/windows_amd64/allmark.exe -OutFile allmark.exe
```

All binaries at [allmark.io](https://allmark.io) are up-to-date builds of the **master**-branch.

If you want to download and install binarier from the **develop**-branch you can go to [develop.allmark.io/bin](https://develop.allmark.io).

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

After a second or so a browser window with this address `http://0.0.0.0:8080` should pop up. Or if it doesn't just type `http://localhost:8080` into your browser:

![Screenshot: Testing the allmark server on the allmark-project directory](files/installation/screenshot-allmark-test-run-on-project-folder.png)

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
	- Repository Cross-Links
17. Different Item Types (Repository, Document, Presentation)
18. Document Meta Data
	⁻ Author
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

I will try to create videos showing you the different features when there is time.

## Demo / Showcase

If you want to see **allmark in action** you can visit my blog [AndyK Docs](https://andykdocs.de/) at [https://andykdocs.de](https://andykdocs.de):

![Animation: Demo of allmark hosting andykdocs.de ](files/demo/allmark-demo-andykdocs.gif)

## Build Status

[![Build Status](https://travis-ci.org/andreaskoch/allmark.png)](https://travis-ci.org/andreaskoch/allmark)

There is also an automated docker build at [registry.hub.docker.com/u/andreaskoch/allmark/](https://registry.hub.docker.com/u/andreaskoch/allmark/) which builds the develop and master branch every time a commit is pushed.

## Dependencies

allmark relies on many great third-party libraries. These are some of them:

- [github.com/bradleypeabody/fulltext](src/github.com/bradleypeabody/fulltext)
- [github.com/gorilla/context](src/github.com/gorilla/context)
- [github.com/gorilla/mux](src/github.com/gorilla/mux)
- [github.com/jbarham/go-cdb](src/github.com/jbarham/go-cdb)
- [github.com/nfnt/resize](src/github.com/nfnt/resize)
- [github.com/russross/blackfriday](src/github.com/russross/blackfriday)
- [github.com/shurcooL/go/github_flavored_markdown/sanitized_anchor_name](src/github.com/shurcooL/go/github_flavored_markdown/sanitized_anchor_name)
- [github.com/skratchdot/open-golang/open](src/github.com/skratchdot/open-golang/open)
- [golang.org/x/net/websocket](src/golang.org/x/net/websocket)
- [github.com/andreaskoch/go-fswatch](src/github.com/andreaskoch/go-fswatch)

These depenendencies are not under allmark copyright/license. See the respective projects for their copyright & licensing details.
The packages are mirrored into allmark for hermetic build reasons and versioning.

To get a full list of all used third-party libraries you execute the make tool with the `-list-dependencies` flag:

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

### Default Theme

- Responsive Design: The default theme always selects the smallest thumbnail size which makes the previewed images look crappy on large screens.

### Windows

- Fileystem links: Serving folders that are fileystem junctions/links is no longer possible with go 1.4 (it did work with go 1.3)

## Roadmap / To Dos

Here are some of the ideas and todos I would like to add in the future. Contributions are welcome!

### Architecture & Features

- Expose the markdown source
- HTTPs support
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
- Make live-reload more intelligent and more efficient
- Support for Folders with multiple Markdown Files

### Theming

- Redesign default theme with Twitter Bootstrap
- Create a theme "loader"
- Infinite Scrolling for latest items
- Improved Image Galleries
