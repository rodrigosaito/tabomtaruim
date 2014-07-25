package main

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
)

type GoodBad struct {
	Line      string `json:"line,omitempty"`
	Imei      string `json:"imei,omitempty"`
	Status    string `json:"status,omitempty"`
	Timestamp uint32 `json:"timestamp,omitempty"`
}

type LineStatus struct {
	Line   string `json:"line,omitempty"`
	Goods  uint32 `json:"goods,omitempty"`
	Bads   uint32 `json:"bads,omitempty"`
	Status string `json:"status,omitempty"`
}

func PostGoodBad(w rest.ResponseWriter, req *rest.Request) {
	goodBad := GoodBad{}
	if err := req.DecodeJsonPayload(&goodBad); err != nil {
		panic(err)
	}

	lineStatus := LineStatus{
		Line:   goodBad.Line,
		Goods:  10,
		Bads:   2,
		Status: "good",
	}

	w.WriteJson(lineStatus)
}

func main() {
	handler := rest.ResourceHandler{}
	handler.SetRoutes(
		&rest.Route{"POST", "/good_bad", PostGoodBad},
	)
	http.ListenAndServe(":8080", &handler)
}
