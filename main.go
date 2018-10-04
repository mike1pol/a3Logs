package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
*/
import "C"

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
	"unsafe"
)

//export RVExtensionVersion
func RVExtensionVersion(output *C.char, outputsize C.size_t) {
	result := C.CString("v0.0.1")
	defer C.free(unsafe.Pointer(result))
	var size = C.strlen(result) + 1
	if size > outputsize {
		size = outputsize
	}
	C.memmove(unsafe.Pointer(output), unsafe.Pointer(result), size)
}

//export RVExtensionArgs
func RVExtensionArgs(output *C.char, outputsize C.size_t, input *C.char, argv **C.char, argc C.int) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	currentDir := filepath.Dir(ex)
	dir := fmt.Sprintf("%s/a3logs", currentDir)
	errMk := os.Mkdir(dir, 0666)
	if errMk != nil {
		panic(errMk)
	}
	filePath := fmt.Sprintf("%s/%s", dir, time.Now().Format("2006-01-02 15:04:05"))
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	var offset = unsafe.Sizeof(uintptr(0))
	type LogInput struct {
		key   string
		value string
	}
	var out []LogInput
	for index := C.int(0); index < argc; index++ {
		if index%2 == 0 {
			key := C.GoString(*argv)
			argv = (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) + offset))
			value := C.GoString(*argv)
			arg := LogInput{
				key:   key,
				value: value,
			}
			out = append(out, arg)
		}
	}
	cmd := C.GoString(input)
	log.Printf("Action: %s, arguments: %s", cmd, out)
	temp := fmt.Sprintf("Action: %s params: %s!", cmd, out)
	result := C.CString(temp)
	defer C.free(unsafe.Pointer(result))
	var size = C.strlen(result) + 1
	if size > outputsize {
		size = outputsize
	}
	C.memmove(unsafe.Pointer(output), unsafe.Pointer(result), size)
}

//export RVExtension
func RVExtension(output *C.char, outputsize C.size_t, input *C.char) {
	temp := fmt.Sprintf("a3-logs ready %s!", C.GoString(input))
	result := C.CString(temp)
	defer C.free(unsafe.Pointer(result))
	var size = C.strlen(result) + 1
	if size > outputsize {
		size = outputsize
	}
	C.memmove(unsafe.Pointer(output), unsafe.Pointer(result), size)
}

func main() {}
