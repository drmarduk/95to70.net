package main

import (
	"database/sql"
	"log"
	"time"
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

	var tmp string
	row := db.QueryRow("select id, datetime(created), value from trend order by created desc limit 0, 1")
	err = row.Scan(&r.ID, &tmp, &r.Value)
	r.Created, err = time.Parse("2006-01-02 15:04:05.999999999", tmp)
	if err != nil {
		log.Printf("error while parsing time from sql column: %s: %v\n", tmp, err)
	}

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
		var x string
		var r Record
		err = rows.Scan(&r.ID, &x, &r.Value)
		if err != nil {
			log.Printf("error while query last items: %s\n", err.Error())
			continue
		}
		r.Created, err = time.Parse("2006-01-02 15:04:05.999999999-07:00", x)
		if err != nil {
			log.Printf("error while parsing time from sql column: %s: %v\n", x, err)
			continue
		}
		result = append(result, r)
	}
	return result, nil
}

// Min returns the minimum value in the trend table
func (w *WeightAPI) Min() (r Record, err error) {
	return w.singleRow("select id, min(created), value from trend")
}

// Max returns the maximum value in the trend table
func (w *WeightAPI) Max() (Record, error) {
	return w.singleRow("select id, max(created), value from trend")
}

func (w *WeightAPI) singleRow(query string, args ...interface{}) (r Record, err error) {
	db, err := sql.Open(w.driver, w.db)
	if err != nil {
		return r, err
	}
	defer db.Close()

	row := db.QueryRow(query, args...)
	var tmp string
	err = row.Scan(&r.ID, &tmp, &r.Value)
	if err != nil {
		return r, err
	}

	r.Created, err = time.Parse("2006-01-02 15:04:05.999999999-07:00", tmp)
	if err != nil {
		log.Printf("error while parsing sql time column: %s: %v\n", tmp, err)
	}
	return r, nil
}
