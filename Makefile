.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/nunu .

.PHONY: build-darwin
build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/nunu .

.PHONY: build-windows
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/nunu.exe .

