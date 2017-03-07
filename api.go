package main

import (
	"database/sql"
	"log"
)

// WeightAPI is the backend api
type WeightAPI struct {
	db     string
	driver string
}

// NewWeightAPI returns a new instance
func NewWeightAPI(driver, db string) *WeightAPI {
	return &WeightAPI{driver: driver, db: db}
}

// Current returns the last inserted value
func (w *WeightAPI) Current() (Record, error) {
	var r Record
	db, err := sql.Open(w.driver, w.db)
	if err != nil {
		return r, err
	}
	defer db.Close()

	row := db.QueryRow("select id, datetime(created), value from trend order by created desc limit 0, 1")
	err = row.Scan(&r.ID, &r.Created, &r.Value)

	log.Println("Current: ", r.String())
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
	log.Println("Add: ", r.String())
	return nil
}

// Last returns the latest n Elements in the thred table
func (w *WeightAPI) Last(n int) ([]Record, error) {
	db, err := sql.Open(w.driver, w.db)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("select id, datetime(created), value  from trend order by created desc limit 0, ?", n)
	if err != nil {
		return nil, err
	}

	var result []Record
	for rows.Next() {
		var r Record
		err = rows.Scan(&r.ID, &r.Created, &r.Value)
		if err != nil {
			log.Printf("error while query last items: %s\n", err.Error())
			continue
		}
		result = append(result, r)
	}
	return result, nil
}

// Min returns the minimum value in the trend table
func (w *WeightAPI) Min() (r Record, err error) {
	return Record{ID: 23, Value: 23.0}, nil
}

// Max returns the maximum value in the trend table
func (w *WeightAPI) Max() (Record, error) {
	return Record{ID: 42, Value: 42.0}, nil
}
