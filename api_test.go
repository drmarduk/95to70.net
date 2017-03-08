package main

import "testing"

func TestNewWeightAPI(t *testing.T) {
	driver := "sqlite3"
	db := ":memory:"
	api := NewWeightAPI(driver, db)

	if api.driver != driver {
		t.Errorf("NewWeightAPI: got %s, expected %s\n", api.driver, driver)
	}
	if api.db != db {
		t.Errorf("NewWeightAPI: got %s, expected %s\n", api.db, db)
	}
}

func TestCurrent(t *testing.T) {

	api := NewWeightAPI("sqlite3", ":memory:")

}

/*
CREATE TABLE `trend` (
	`id`	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	`created`	INTEGER NOT NULL,
	`value`	REAL NOT NULL
);

*/
