# allmark - the tiny markdown server

allmark is a lightweight markdown web server for Linux, Mac OS and Windows written in go.

## Build Status

[![Build Status](https://travis-ci.org/andreaskoch/allmark.png)](https://travis-ci.org/andreaskoch/allmark)

## Dependencies

- [Blackfriday: a markdown processor for Go](https://github.com/russross/blackfriday)
- [go-fswatch: a library for monitoring file system changes](https://github.com/andreaskoch/go-fswatch)

## Roadmap / To Dos

- Run on Raspberry Pi / WDLXTV ("Host your blog from your home")
    - store images and thumbnails on amazon s3
    - can be run with very little bandwidth
    - DDNS support
- Data Access
    - Dropbox support
    - SMTP message posting
    - Repository Replication?
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
    - tag links
    - geo locations
- Different file formats
    - json
- Web Server
    - Cache Header Management
- Search
    - Amazon Cloud Search?
    - Lucene?