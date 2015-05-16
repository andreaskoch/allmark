# allmark - the markdown server

allmark is a standalone markdown web server for Linux, BSD, Solaris Mac OS and Windows written in go.

![allmark logo (128x128px)](files/design/logo/PNG8/allmark-logo-128x128.png)

## Build Status

[![Build Status](https://travis-ci.org/andreaskoch/allmark.png)](https://travis-ci.org/andreaskoch/allmark)

## Installation

If you have [go](https://golang.org/dl/) (≥ 1.3) installed all you have to do is

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

## Demo / Showcase

If you want to see **allmark in action** you can visit my blog [AndyK Docs](https://andykdocs.de/) at [https://andykdocs.de](https://andykdocs.de):

![Animation: Demo of allmark hosting andykdocs.de ](files/demo/allmark-demo-andykdocs.gif)

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

### Theming

- Redesign default theme with Twitter Bootstrap
- Create a theme "loader"
- Infinite Scrolling for latest items
- Improved Image Galleries
