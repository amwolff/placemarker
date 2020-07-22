package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/amwolff/placemarker"
	"github.com/twpayne/go-kml"
)

type locationInfo struct {
	Name string  `json:"name"`
	Lon  float64 `json:"lon"`
	Lat  float64 `json:"lat"`
	Alt  float64 `json:"alt"`
}

type stateful struct {
	k    *kml.CompoundElement
	mu   sync.Mutex
	path string
}

func (s *stateful) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var info locationInfo
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		log.Printf("ServeHTTP: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.k = placemarker.AddPoint(s.k, info.Name, info.Lon, info.Lat, info.Alt)

	if err := placemarker.WriteKML(s.k, s.path); err != nil {
		log.Printf("ServeHTTP: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Fprintln(w, http.StatusText(http.StatusOK))
}

func getMain() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func main() {
	s := &stateful{
		path: "database.kml",
	}

	http.Handle("/", getMain())
	http.Handle("/insert", s)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Printf("ListenAndServe: %s", err)
	}
}
