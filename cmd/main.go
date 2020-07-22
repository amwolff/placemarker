package main

import (
	"encoding/json"
	"net/http"

	"github.com/twpayne/go-kml"
)

type locationInfo struct {
	Name string  `json:"name"`
	Desc string  `json:"desc"`
	Alt  float64 `json:"alt"`
	Lon  float64 `json:"lon"`
	Lat  float64 `json:"lat"`
}

func getInsertPlacemark() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := r.GetBody()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var location locationInfo
		if err := json.NewDecoder(b).Decode(&location); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func main() {
	var mainKML *kml.CompoundElement
}
