package main

import (
	"reflect"
	"testing"
)

const testJSONPath = "./test_fixtures/test.json"

func TestNewWhiteList(t *testing.T) {
	wl, err := NewWhitelist(testJSONPath)
	if err != nil {
		t.Error("Whitelist not created: ", err)
	}
	expectedWhiteList := Whitelist{
		Databases{
			"common_development": Tables{
				"account_shard_mappings": Colnames{"id", "account_id", "shard_id", "name", "created_at", "updated_at"},
				"admin_bonus_points":     Colnames{"id", "admin_id", "reward_tracker_id", "reason", "created_at"},
			},
			"another_db": Tables{
				"table1": Colnames{"col1", "col2"},
			},
		},
	}

	if !reflect.DeepEqual(*wl, expectedWhiteList) {
		t.Error("Parsing went wrong")
	}
}
