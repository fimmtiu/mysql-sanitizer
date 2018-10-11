package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Colnames []string

type Tables map[string]Colnames

type Databases map[string]Tables

type Whitelist struct {
	Databases Databases
}

func NewWhitelist(path string) (*Whitelist, error) {
	wl := &Whitelist{}
	db := Databases{}
	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = json.Unmarshal(byteValue, &db)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	wl.Databases = db
	return wl, nil

}
