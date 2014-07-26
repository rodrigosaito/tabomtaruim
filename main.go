package main

import (
	"net/http"
	"os"

	"github.com/rodrigosaito/tabomtaruim/models"

	"github.com/ant0ine/go-json-rest/rest"
	"gopkg.in/mgo.v2"
)

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
	goodBad := models.GoodBad{}
	if err := req.DecodeJsonPayload(&goodBad); err != nil {
		panic(err)
	}

	goodBad.Save(api.Db)

	lineStatus := models.LineStatus{
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

	mongoUrl := os.Getenv("MONGOHQ_URL")
	if mongoUrl == "" {
		mongoUrl = "localhost"
	}

	mongoDbName := os.Getenv("MONGO_DATABASE_NAME")
	if mongoDbName == "" {
		mongoDbName = "good_bad_dev"
	}

	handler := rest.ResourceHandler{}

	api := &GoodBadApi{
		MongoUrl: mongoUrl,
		DbName:   mongoDbName,
	}
	api.Init()

	handler.SetRoutes(
		rest.RouteObjectMethod("POST", "/good_bad", api, "PostGoodBad"),
	)
	http.ListenAndServe(":"+port, &handler)
}
