# allmark - the tiny markdown server

allmark is a lightweight markdown web server for Linux, Mac OS and Windows written in go.

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

and [github.com/laher/goxc](src/github.com/laher/goxc) for cross-compilation.

These are not under allmark copyright/license. See the respective projects for their copyright & licensing details.
These are mirrored into allmark for hermetic build reasons and versioning.

To get a full list of all used third-party libraries you execute the make tool with the `-dependencies` flag:

```bash
go run make.go -dependencies
```

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
