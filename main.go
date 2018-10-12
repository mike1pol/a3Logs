package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
*/
import "C"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unsafe"
)

func send(output *C.char, outputsize C.size_t, data *C.char) {
	defer C.free(unsafe.Pointer(data))
	var size = C.strlen(data) + 1
	if size > outputsize {
		size = outputsize
	}
	C.memmove(unsafe.Pointer(output), unsafe.Pointer(data), size)
}

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
	if config == nil {
		errConfig := getConfig()
		if errConfig != nil {
			error := C.CString(fmt.Sprintf("Error open config file: %s", errConfig))
			send(output, outputsize, error)
			return
		}
	}
	if templates == nil {
		getTemplate()
	}

	if db == nil {
		errDB := getDB()
		if errDB != nil {
			error := C.CString(fmt.Sprintf("Error connecting to database: %s", errDB))
			send(output, outputsize, error)
			return
		}
	}

	var offset = unsafe.Sizeof(uintptr(0))
	out := make(map[string]string)
	out["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	var arr []string
	for index := C.int(0); index < argc; index++ {
		arr = append(arr, strings.Replace(C.GoString(*argv), "\"", "", -1))
		argv = (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) + offset))
	}
	for index := 0; index < len(arr); index++ {
		if index%2 == 0 {
			key := arr[index]
			value := arr[index+1]
			out[key] = value
		}
	}
	cmd := C.GoString(input)
	tmpl, ok := templates[cmd]
	if !ok {
		error := C.CString(fmt.Sprintf("Error template `%s` not found", cmd))
		send(output, outputsize, error)
		return
	}

	type templType struct {
		Map  map[string]string
		JSON string
	}
	jData, _ := json.Marshal(&out)
	jsonData := strings.Replace(string(jData), `"`, `\"`, -1)
	data := templType{
		Map:  out,
		JSON: jsonData,
	}

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, data)
	res, err := db.Exec(tpl.String())

	var outMSG string

	if err != nil {
		outMSG = fmt.Sprintf("Error insert log: %s, sql: %s", err, tpl.String())
	} else {
		id, _ := res.LastInsertId()
		outMSG = fmt.Sprintf("Log for player: %s id: %d!", cmd, id)
	}
	result := C.CString(outMSG)
	send(output, outputsize, result)
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
