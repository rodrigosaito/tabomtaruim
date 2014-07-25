package main

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
)

type GoodBad struct {
	Imei      string `json:"imei,omitempty"`
	Status    string `json:"status,omitempty"`
	Timestamp uint32 `json:"timestamp,omitempty"`
}

func PostGoodBad(w rest.ResponseWriter, req *rest.Request) {
	goodBad := GoodBad{}
	if err := req.DecodeJsonPayload(&goodBad); err != nil {
		panic(err)
	}

	w.WriteJson(goodBad)
}

func main() {
	handler := rest.ResourceHandler{}
	handler.SetRoutes(
		&rest.Route{"POST", "/good_bad", PostGoodBad},
	)
	http.ListenAndServe(":8080", &handler)
}
