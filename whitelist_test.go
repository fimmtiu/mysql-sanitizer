package main

import (
	"reflect"
	"testing"
)

const testJSONPath = "./test_fixtures/test.json"

func TestNewWhitelist(t *testing.T) {
	wl, err := NewWhitelist(testJSONPath)
	if err != nil {
		t.Error("Whitelist not created: ", err)
	}
	expectedWhiteList := Whitelist{
		Databases{
			"some_db": Tables{
				"table1": Colnames{"id", "foo_id", "bar_id", "name", "created_at", "updated_at"},
				"table2": Colnames{"id", "honk", "bonk", "created_at"},
			},
			"another_db": Tables{
				"table1": Colnames{"col1", "col2"},
			},
		},
	}

	if !reflect.DeepEqual(wl, expectedWhiteList) {
		t.Error("Parsing went wrong")
	}
}

func TestWhitelistIsColumnPresent(t *testing.T) {
	wl, err := NewWhitelist(testJSONPath)
	if err != nil {
		t.Error("Whitelist not created: ", err)
	}

	if !wl.IsColumnPresent("some_db", "table2", "bonk") {
		t.Error("Couldn't find a column that should exist")
	}

	if wl.IsColumnPresent("some_db", "table2", "arghablargh") {
		t.Error("Found a nonexistent column")
	}

	if wl.IsColumnPresent("some_db", "table31337", "bonk") {
		t.Error("Found a column in a nonexistent table")
	}

	if wl.IsColumnPresent("gooooooooats", "table2", "bonk") {
		t.Error("Found a column in a nonexistent database")
	}
}
