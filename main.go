package main

import (
	"html/template"
	"io"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
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

func renderIndex(w http.ResponseWriter) (err error) {
	t := template.New("index")
	t, err = t.ParseFiles("./html/index.html")
	print("asdf")
	if err != nil {
		log.Printf("renderIndex: %v\n", err)
		return err
	}
	return t.Execute(w, nil)
}

// handler stuff
func indexHandler(w http.ResponseWriter, r *http.Request) error {
	return renderIndex(w)
}

func addHandler(w http.ResponseWriter, r *http.Request) error {
	record, err := ParseRecord(r.FormValue("value"))
	if err != nil {
		return err
	}

	api := NewWeightAPI()
	err = api.Add(record)
	if err != nil {
		return err
	}

	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}
