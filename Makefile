build-osx: btc.go
	GOOS=darwin GOARCH=amd64 go build main.go

build-win-64: btc.go
	GOOS=windows GOARCH=amd64 go build main.go