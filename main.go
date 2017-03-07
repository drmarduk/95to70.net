package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var _driver = "sqlite3"
var _db = "test.db"

func main() {

	fs := http.FileServer(http.Dir("./html"))
	http.Handle("/favico.ico", fs)

	http.HandleFunc("/", errorHandler(indexHandler))
	http.HandleFunc("/api/add", errorHandler(addHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// render stuff
func error502(w http.ResponseWriter, err error) {
	io.WriteString(w, err.Error())
}

func errorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("handling %q: %v\n", r.RequestURI, err)
		}
	}
}

func renderIndex(w http.ResponseWriter, s Summary) (err error) {
	t := template.New("index")
	t, err = t.ParseFiles("./html/index.html")

	if err != nil {
		log.Printf("renderIndex: %v\n", err)
		return err
	}
	return t.Execute(w, nil)
}

// handler stuff
func indexHandler(w http.ResponseWriter, r *http.Request) (err error) {
	s := Summary{}

	api := NewWeightAPI(_driver, _db)

	s.Current, err = api.Current()
	if err != nil {
		s.Err = append(s.Err, err.Error())
	}

	s.Min, err = api.Min()
	if err != nil {
		s.Err = append(s.Err, err.Error())
	}

	s.Max, err = api.Min()
	if err != nil {
		s.Err = append(s.Err, err.Error())
	}

	s.LastMonth, err = api.Last(30)
	if err != nil {
		s.Err = append(s.Err, err.Error())
	}

	for _, v := range s.Err {
		log.Println(v)
	}
	fmt.Printf("Summary: %v\n\n", s)
	return renderIndex(w, s)
}

func addHandler(w http.ResponseWriter, r *http.Request) error {
	// TODO: add authorization
	record, err := ParseRecord(r.FormValue("value"))
	if err != nil {
		return err
	}

	api := NewWeightAPI(_driver, _db)

	err = api.Add(record)
	if err != nil {
		return err
	}

	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}

// Summary holds all info for the frontpage
type Summary struct {
	Err       []string
	Current   Record
	Min       Record
	Max       Record
	LastMonth []Record
}
