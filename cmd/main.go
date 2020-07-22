package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

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

	log.Printf("Got: %v", info)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.k = placemarker.AddPoint(s.k, info.Name, info.Lon, info.Lat, info.Alt)

	p := filepath.Join(s.path, fmt.Sprintf("%s.kml", strconv.FormatInt(time.Now().UnixNano(), 10)))

	if err := placemarker.WriteKML(s.k, p); err != nil {
		log.Printf("ServeHTTP: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, http.StatusText(http.StatusOK))
}

func getMain(index []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(index))
	}
}

func main() {
	confPath := os.Getenv("CONFDIR")
	dataPath := os.Getenv("DATADIR")

	s := &stateful{
		path: dataPath,
	}

	index, err := ioutil.ReadFile(filepath.Join(confPath, "index.html"))
	if err != nil {
		log.Panicf("ReadFile: %s", err)
	}

	http.Handle("/", getMain(index))
	http.Handle("/insert", s)

	cert := filepath.Join(confPath, "cert.pem")
	priv := filepath.Join(confPath, "privkey.pem")

	if err := http.ListenAndServeTLS(":https", cert, priv, nil); err != nil {
		log.Printf("ListenAndServeTLS: %s", err)
	}
}
