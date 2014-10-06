# allmark - the tiny markdown server

allmark is a lightweight markdown web server for Linux, Mac OS and Windows written in go.

## Build Status

[![Build Status](https://travis-ci.org/andreaskoch/allmark.png)](https://travis-ci.org/andreaskoch/allmark)

## Dependencies

- [Blackfriday: a markdown processor for Go](https://github.com/russross/blackfriday)
- [bradleypeabody/fulltext: a Pure-Go full text indexer and search library](https://github.com/bradleypeabody/fulltext)
- [Gorilla web toolkit: mux for request routing](http://www.gorillatoolkit.org/pkg/mux)

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