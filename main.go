package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

}

// Record holds a current weight
type Record struct {
	ID      int
	Created time.Time
	Value   float32
}

// WeightAPI is the backend api
type WeightAPI struct {
	db     string
	driver string
}

// NewWeightAPI returns a new instance
func NewWeightAPI(db string) *WeightAPI {
	return &WeightAPI{db: db, driver: "sqlite3"}
}

// Current returns the last inserted value
func (w *WeightAPI) Current() (*Record, error) {
	db, err := sql.Open(w.driver, w.db)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	row := db.QueryRow("select * from trend order by creadted desc limit 0, 1")

	result := &Record{}
	err = row.Scan(&result.ID, &result.Created, &result.Value)

	return result, err
}
