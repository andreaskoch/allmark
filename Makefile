build:
	go build -o bin/files/ ./cmd/...

install:
	go build -o bin/files/ ./cmd/...

test:
	go test ./common/... ./cmd/... ./dataaccess/... ./model/... ./services/... ./web/...

crosscompile:
	GOOS=linux GOARCH=amd64 go build -o bin/files/allmark_linux_amd64 ./cmd/allmark
	GOOS=linux GOARCH=arm GOARM=5 go build -o bin/files/allmark_linux_arm_5 ./cmd/allmark
	GOOS=linux GOARCH=arm GOARM=6 go build -o bin/files/allmark_linux_arm_6 ./cmd/allmark
	GOOS=linux GOARCH=arm GOARM=7 go build -o bin/files/allmark_linux_arm_7 ./cmd/allmark
	GOOS=darwin GOARCH=amd64 go build -o bin/files/allmark_darwin_amd64 ./cmd/allmark
	GOOS=windows GOARCH=amd64 go build -o bin/files/allmark_windows_amd64 ./cmd/allmark
