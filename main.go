package main

import (
	"fmt"
	"net/http"

	"github.com/MShoaei/FlyMode/flight/eligasht"
	"github.com/kataras/muxie"
)

var flightWebsites map[string]func(http.ResponseWriter, *http.Request)

func init() {
	flightWebsites = make(map[string]func(http.ResponseWriter, *http.Request), 20)
	flightWebsites["eligasht"] = eligasht.FindFlights
}

func flight(w http.ResponseWriter, r *http.Request) {
	website := muxie.GetParam(w, "name")
	handler := flightWebsites[website]
	if handler == nil {
		fmt.Fprintf(w, "%s is not implemented", website)
		return
	}
	handler(w, r)
}

func main() {
	mux := muxie.NewMux()
	mux.PathCorrection = true

	mux.HandleFunc("/flight/:name", flight)
	http.ListenAndServe(":8085", mux)
}
