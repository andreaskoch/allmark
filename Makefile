build:
	go build -o bin/files/allmark ./cli

install:
	go build -o bin/files/allmark ./cli

test:
	go test ./cli ./common/... ./dataaccess/... ./model/... ./services/... ./web/...

crosscompile:
	GOOS=linux GOARCH=amd64 go build -o bin/files/allmark_linux_amd64 ./cli
	GOOS=linux GOARCH=arm GOARM=5 go build -o bin/files/allmark_linux_arm_5 ./cli
	GOOS=linux GOARCH=arm GOARM=6 go build -o bin/files/allmark_linux_arm_6 ./cli
	GOOS=linux GOARCH=arm GOARM=7 go build -o bin/files/allmark_linux_arm_7 ./cli
	GOOS=darwin GOARCH=amd64 go build -o bin/files/allmark_darwin_amd64 ./cli
	GOOS=windows GOARCH=amd64 go build -o bin/files/allmark_windows_amd64 ./cli
