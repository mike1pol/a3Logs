package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/ini.v1"
)

var templates map[string]*template.Template
var config *ini.File
var db *sql.DB

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

func getConfig() error {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	cFile := fmt.Sprintf("%s/@a3Logs/%s", dir, "config.ini")
	cfg, err := ini.Load(cFile)
	config = cfg
	if err != nil {
		return fmt.Errorf("%s - %s", cFile, err)
	}
	return nil
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
