package main

import (
	"database/sql"
	"testing"
)

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
	// api := NewWeightAPI("sqlite3", ":memory:")

	db, err := sql.Open("sqlite3", ":memory:?parseTime=true")
	if err != nil {
		t.Errorf("Failed to open database: %v\n", err)
	}
	defer db.Close()

	var d string
	err = db.QueryRow("select cast(datetime('now') as text)").Scan(&d)
	if err != nil {
		t.Errorf("Failed to scan datetime: %v\n", err)
	}
}

/*
CREATE TABLE `trend` (
	`id`	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	`created`	INTEGER NOT NULL,
	`value`	REAL NOT NULL
);

*/
// 92; 91,6; 91; 93; 91,5
