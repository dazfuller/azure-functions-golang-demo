build:
	go build -o app cmd/main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -o app cmd/main.go

run: build
	func start
