# The normal way to build allmark is just "go run make.go", which
# works everywhere, even on systems without Make.

all:
	go run make.go -install

fmt:
	go run make.go -fmt