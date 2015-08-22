# Source

All of allmarks' sources

## Dependencies

allmark relies on a number of third-party libraries:

- [github.com/bradleypeabody/fulltext](github.com/bradleypeabody/fulltext)
- [github.com/gorilla/context](github.com/gorilla/context)
- [github.com/gorilla/mux](github.com/gorilla/mux)
- [github.com/gorilla/handlers](github.com/gorilla/handlers)
- [github.com/jbarham/go-cdb](github.com/jbarham/go-cdb)
- [github.com/nfnt/resize](github.com/nfnt/resize)
- [github.com/russross/blackfriday](github.com/russross/blackfriday)
- [github.com/shurcooL/go/github_flavored_markdown/sanitized_anchor_name](github.com/shurcooL/go/github_flavored_markdown/sanitized_anchor_name)
- [github.com/skratchdot/open-golang/open](github.com/skratchdot/open-golang/open)
- [golang.org/x/net/websocket](golang.org/x/net/websocket)
- [github.com/andreaskoch/go-fswatch](github.com/andreaskoch/go-fswatch)
- [github.com/abbot/go-http-auth](github.com/abbot/go-http-auth)
- [github.com/spf13/afero](github.com/spf13/afero)
- [github.com/kyokomi/emoji](github.com/kyokomi/emoji)

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

---

created at: 2015-08-03
modified at: 2015-08-03
author: Andreas Koch
tags: Source Code
alias: src, source, code
