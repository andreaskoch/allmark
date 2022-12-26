module github.com/andreaskoch/allmark

go 1.19

// remove when andreaskoch/fulltext#1 closes
replace github.com/andreaskoch/fulltext => github.com/urandom2/fulltext v0.0.0-20221226014327-4b3d48bf6613

require (
	github.com/abbot/go-http-auth v0.4.0
	github.com/andreaskoch/fulltext v0.0.0-00010101000000-000000000000
	github.com/andreaskoch/go-fswatch v1.0.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/kyokomi/emoji v1.5.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/russross/blackfriday v1.6.0
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966
	github.com/spf13/afero v1.9.3
	golang.org/x/net v0.4.0
)

require (
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/jbarham/cdb v0.0.0-20200301055225-9d6f6caadef0 // indirect
	golang.org/x/crypto v0.0.0-20211108221036-ceb1ce70b4fa // indirect
	golang.org/x/text v0.5.0 // indirect
)
