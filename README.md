# a3Logs to DB [![Build Status](https://travis-ci.org/mike1pol/a3Logs.svg?branch=master)](https://travis-ci.org/mike1pol/a3Logs)
## Requirements & Build
* Requirements: Golang (https://golang.org), Dep (https://github.com/golang/dep)
* For windows https://sourceforge.net/projects/tdm-gcc/

### Build
1. install dependency `dep ensure`

2. build your extension with this command line:
windows: `go build -o a3Logs_x64.dll -buildmode=c-shared .`
linux: `go build -o a3Logs_x64.so  -buildmode=c-shared .`
