package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
*/
import "C"

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
	"unsafe"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/ini.v1"
)

var templates map[string]*template.Template
var db *sql.DB
var config *ini.File

func getConfig() error {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	cFile := fmt.Sprintf("%s/@a3Logs/%s", dir, "config.ini")
	cfg, err := ini.Load(cFile)
	config = cfg
	return fmt.Errorf("%s - %s", cFile, err)
}

func getTemplate() {
	tmpls := config.Section("templates").KeyStrings()
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	for _, t := range tmpls {
		mt := template.New(t)
		mt, err := mt.Parse(config.Section("templates").Key(t).String())
		if err == nil {
			templates[t] = mt
		}

	}
}

func getDB() error {
	dbURL := config.Section("").Key("db").String()
	if len(dbURL) == 0 {
		return errors.New("dbUrl not set in config file")
	}
	d, err := sql.Open("mysql", dbURL)
	if err != nil {
		return err
	}
	db = d
	return nil
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
			temp := fmt.Sprintf("Error open config file: %s", errConfig)
			result := C.CString(temp)
			defer C.free(unsafe.Pointer(result))
			var size = C.strlen(result) + 1
			if size > outputsize {
				size = outputsize
			}
			C.memmove(unsafe.Pointer(output), unsafe.Pointer(result), size)
			return
		}
	}
	if templates == nil {
		getTemplate()
	}

	if db == nil {
		errDB := getDB()
		if errDB != nil {
			temp := fmt.Sprintf("Error connecting to database: %s", errDB)
			result := C.CString(temp)
			defer C.free(unsafe.Pointer(result))
			var size = C.strlen(result) + 1
			if size > outputsize {
				size = outputsize
			}
			C.memmove(unsafe.Pointer(output), unsafe.Pointer(result), size)
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
		temp := fmt.Sprintf("Error template `%s` not found", cmd)
		result := C.CString(temp)
		defer C.free(unsafe.Pointer(result))
		var size = C.strlen(result) + 1
		if size > outputsize {
			size = outputsize
		}
		C.memmove(unsafe.Pointer(output), unsafe.Pointer(result), size)
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

	d1 := []byte(fmt.Sprintf("arr: %s\r\njson: %s\r\nmap: %s\r\nsql: %s", arr, jsonData, out, tpl.String()))
	ioutil.WriteFile("sql.log", d1, 0644)

	var outMSG string

	if err != nil {
		outMSG = fmt.Sprintf("Error insert log: %s, sql: %s", err, tpl.String())
	} else {
		id, _ := res.LastInsertId()
		outMSG = fmt.Sprintf("Log for player: %s id: %d!", cmd, id)
	}
	result := C.CString(outMSG)
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
