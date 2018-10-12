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

func NewWhitelist(path string) (Whitelist, error) {
	wl := Whitelist{}
	db := Databases{}
	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
		return wl, err
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return wl, err
	}

	err = json.Unmarshal(byteValue, &db)
	if err != nil {
		fmt.Println(err)
		return wl, err
	}
	wl.Databases = db
	return wl, nil
}

// TODO: If we made the last component a hash instead of an array, the JSON
// would be uglier but lookups would be cheaper.
func (wl Whitelist) IsColumnPresent(database string, table string, colname string) bool {
	if _, ok := wl.Databases[database]; ok {
		if _, ok := wl.Databases[database][table]; ok {
			for _, name := range wl.Databases[database][table] {
				if colname == name {
					return true
				}
			}
		}
	}
	return false
}
