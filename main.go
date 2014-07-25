package main

import (
	// "log"
	"net/http"
	"os"

	"github.com/ant0ine/go-json-rest/rest"
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
)

type GoodBad struct {
	Line      string `json:"line,omitempty"`
	Imei      string `json:"imei,omitempty"`
	Status    string `json:"status,omitempty"`
	Timestamp uint32 `json:"timestamp,omitempty"`
}

func (gb *GoodBad) Collection(db *mgo.Database) *mgo.Collection {
	return db.C("good_bad")
}

func (gb *GoodBad) Save(db *mgo.Database) {
	if err := gb.Collection(db).Insert(gb); err != nil {
		panic(err)
	}
}

type LineStatus struct {
	Line   string `json:"line,omitempty"`
	Goods  uint32 `json:"goods,omitempty"`
	Bads   uint32 `json:"bads,omitempty"`
	Status string `json:"status,omitempty"`
}

type GoodBadApi struct {
	MongoUrl string
	DbName   string
	// Session  *mgo.Session
	Db *mgo.Database
}

func (api *GoodBadApi) Init() {
	session, err := mgo.Dial(api.MongoUrl)
	if err != nil {
		panic(err)
	}

	api.Db = session.DB(api.DbName)
}

func (api *GoodBadApi) PostGoodBad(w rest.ResponseWriter, req *rest.Request) {
	goodBad := GoodBad{}
	if err := req.DecodeJsonPayload(&goodBad); err != nil {
		panic(err)
	}

	goodBad.Save(api.Db)

	lineStatus := LineStatus{
		Line:   goodBad.Line,
		Goods:  10,
		Bads:   2,
		Status: "good",
	}

	w.WriteJson(lineStatus)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	handler := rest.ResourceHandler{}

	api := &GoodBadApi{
		MongoUrl: "localhost",
		DbName:   "good_bad_dev",
	}
	api.Init()

	handler.SetRoutes(
		rest.RouteObjectMethod("POST", "/good_bad", api, "PostGoodBad"),
	)
	http.ListenAndServe(":"+port, &handler)
}
