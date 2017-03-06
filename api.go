package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

// Record holds a current weight
type Record struct {
	ID      int
	Created time.Time
	Value   float32
}

func (r *Record) String() string {
	return fmt.Sprintf("%d - %.2f (%s)", r.ID, r.Value, r.Created.Format(time.UnixDate))
}

// ParseRecord returns a fully parsed record based on the input string
// and adds the created member
func ParseRecord(value string) (r Record, err error) {
	if value == "" {
		return r, errors.New("empty value")
	}
	x, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return r, err
	}

	r.Created = time.Now()
	r.Value = float32(x)
	return
}

// WeightAPI is the backend api
type WeightAPI struct {
	db     string
	driver string
}

// NewWeightAPI returns a new instance
func NewWeightAPI() *WeightAPI {
	return &WeightAPI{db: "test.db", driver: "sqlite3"}
}

// Current returns the last inserted value
func (w *WeightAPI) Current() (Record, error) {
	var r Record
	db, err := sql.Open(w.driver, w.db)
	if err != nil {
		return r, err
	}
	defer db.Close()

	row := db.QueryRow("select * from trend order by created desc limit 0, 1")

	err = row.Scan(&r.ID, &r.Created, &r.Value)

	log.Println("Current: ", r)
	return r, err
}

// Add inserts a new record in the database
// may be
func (w *WeightAPI) Add(r Record) error {
	db, err := sql.Open(w.driver, w.db)
	if err != nil {
		return err
	}
	defer db.Close()

	x, err := db.Exec("insert into trend(id, created, value) values(null, ?, ?)", r.Created, r.Value)
	if err != nil {
		return err
	}

	y, _ := x.LastInsertId()
	r.ID = int(y)
	log.Println("Add: ", r)
	return nil
}
