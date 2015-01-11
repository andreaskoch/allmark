# allmark - the tiny markdown server

allmark is a lightweight markdown web server for Linux, BSD, Solaris Mac OS and Windows written in go.

## Build Status

[![Build Status](https://travis-ci.org/andreaskoch/allmark.png)](https://travis-ci.org/andreaskoch/allmark)

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

These are not under allmark copyright/license. See the respective projects for their copyright & licensing details.
These are mirrored into allmark for hermetic build reasons and versioning.

To get a full list of all used third-party libraries you execute the make tool with the `-dependencies` flag:

```bash
go run make.go -dependencies
```

## Cross-Compilation

If you want to cross-compile allmark for different platforms and architectures you can do so by using the `-crosscompile` flag for the make script (if you have [docker](https://www.docker.com) >= 1.4 installed):

```bash
go run make.go -crosscompile
```

This command will launch a [docker container with go 1.4](https://registry.hub.docker.com/u/library/golang/) in it that is prepared for cross-compilation and build allmark for you. The output will be available in the `bin` folder of this project:

```
bin/
├── allmark
├── darwin_386
│   └── allmark
├── darwin_amd64
│   └── allmark
├── dragonfly_386
│   └── allmark
├── dragonfly_amd64
│   └── allmark
├── freebsd_386
│   └── allmark
├── freebsd_amd64
│   └── allmark
├── freebsd_arm
│   └── allmark
├── linux_386
│   └── allmark
├── linux_arm
│   └── allmark
├── nacl_386
│   └── allmark
├── nacl_amd64p32
│   └── allmark
├── nacl_arm
│   └── allmark
├── netbsd_386
│   └── allmark
├── netbsd_amd64
│   └── allmark
├── netbsd_arm
│   └── allmark
├── openbsd_386
│   └── allmark
├── openbsd_amd64
│   └── allmark
├── solaris_amd64
│   └── allmark
├── windows_386
│   └── allmark.exe
├── windows_amd64
│   └── allmark.exe
├── ...
├── README.md
└── src

```

If you don't have docker or don't want to install it you can use [goxc](https://github.com/laher/goxc) to cross-compile allmark.

## Roadmap / To Dos

- Expose the markdown source
- Redesign with Twitter Bootstrap
    - Lazy Loading for Images
    - Smaller Footprint -> require js?
- Infinite Scrolling
    - [jQuery Hash Change](http://benalman.com/code/projects/jquery-hashchange/examples/hashchange/)
- Run on Raspberry Pi / WDLXTV ("Host your blog from your home")
    - store images and thumbnails on amazon s3
    - can be run with very little bandwidth
    - DDNS support
- Data Access
    - Dropbox support
    - SMTP message posting
    - Repository Replication?
    - Amazon S3
- allmark swarm
    - Repository sharding
    - load-balancing
- User Management / Access Restrictions
    - User management pages
- Editing
    - sublime snippets
    - sublime theme
    - Examples
- Update content on change
    - auto update for local files
    - javascript path fixer for local files
- Posting comments
- Rendering / Markdown extensions
    - 360° panoramas
    - image galleries (implemented but needs improvement)
    - file lists (implemented but needs improvement)
    - cross references
    - geo locations
- Different file formats
    - json
- Web Server
    - Cache Header Management
    - GZIP Compression
- Search
    - Amazon Cloud Search?
    - Lucene?
