# Building allmark

You can build allmark yourself by executing the make.go script.

If you have [go](HTTPS://golang.org/dl/) (≥ 1.3) installed you can build allmark yourself in two steps:

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

![Screenshot: Testing the allmark server on the allmark-project directory](files/screenshot-allmark-test-run-on-project-folder.png)

## Cross-Compilation

If you want to cross-compile allmark for different platforms (darwin, dragonfly, freebsd, linux, nacl, netbsd, openbsd, windows) and architectures (amd64, arm) you can do so by using the `-crosscompile` or `-crosscompile-with-docker` flags for the make script.

If you have prepared your go environment for cross-compilation (see: [Dave Cheney - An introduction to cross compilation with Go](http://dave.cheney.net/2012/09/08/an-introduction-to-cross-compilation-with-go)) you can use the `-crosscompile` flag:

```bash
go run make.go -crosscompile
```

This command will cross-compile for all platforms and architectures directly on your system.

If you have not prepared your golang installation for cross-compilation you can use the the `-crosscompile-with-docker` flag instead. This command will launch a [docker container with go 1.4](HTTPS://registry.hub.docker.com/u/library/golang/) that is prepared for cross-compilation and will build allmark for all available platforms and architectures inside the docker-container. The output will be available in the `bin` folder of this project:

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

---

created at: 2015-08-05
modified at: 2015-08-22
author: Andreas Koch
tags: Documentation, Build
alias: build, compile
