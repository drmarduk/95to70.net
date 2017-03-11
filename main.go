package main

import (
	"html/template"
	"io"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var _driver = "sqlite3"
var _db = "test.db"

func main() {

	http.HandleFunc("/", errorHandler(indexHandler))
	http.HandleFunc("/api/add", errorHandler(addHandler))
	http.HandleFunc("/favicon.ico", http.NotFound)
	http.Handle("/static/", http.FileServer(http.Dir("./html/")))
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
	s.Derp = "looooool"
	t := template.New("index")
	t, err = t.ParseFiles("./html/index.html")

	if err != nil {
		log.Printf("renderIndex: %v\n", err)
		return err
	}
	return t.Execute(w, &s)
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

	s.Max, err = api.Max()
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
	return renderIndex(w, s)
}

func addHandler(w http.ResponseWriter, r *http.Request) error {
	// TODO: add authorization
	record, err := ParseRecord(r.FormValue("value"))
	if err != nil {
		log.Printf("API ERROR: add: %v\n", err)
		return err
	}

	api := NewWeightAPI(_driver, _db)

	err = api.Add(record)
	if err != nil {
		log.Printf("API ERROR: add: %v\n", err)
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
	Derp      string
}
