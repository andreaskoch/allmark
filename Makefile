# The normal way to build allmark is just "go run make.go", which
# works everywhere, even on systems without Make.

install:
	go run make.go -install

test:

	go run make.go -test

fmt:
	go run make.go -fmt