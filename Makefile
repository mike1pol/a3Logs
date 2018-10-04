all: clean build
clean:
	rm -rf main
deps:
	dep ensure
build: *.go
	env GOOS=linux GOARCH=amd64 go build -o a3Logs_x64.so -buildmode=c-shared .
	env GOOS=linux GOARCH=386 go build -o a3Logs.so -buildmode=c-shared .
	env GOOS=windows GOARCH=amd64 go build -o a3Logs_x64.dll -buildmode=c-shared .
	env GOOS=windows GOARCH=386 go build -o a3Logs.dll -buildmode=c-shared .
