all: clean deps build
clean:
	rm -rf *.dll
	rm -rf *.h
deps:
	dep ensure
build: *.go
	env CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -o a3Logs/a3Logs_x64.dll -buildmode=c-shared .
	env CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386 CGO_ENABLED=1 go build -o a3Logs/a3Logs.dll -buildmode=c-shared .
	env GOOS=linux GOARCH=386 CGO_ENABLED=1 go build -o a3Logs/a3Logs.so -buildmode=c-shared .
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o a3Logs/a3Logs_x64.so -buildmode=c-shared .
