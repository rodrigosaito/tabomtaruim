package main

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/rodrigosaito/tabomtaruim/models"

	"github.com/ant0ine/go-json-rest/rest"
	"gopkg.in/mgo.v2"
)

type GoodBadApi struct {
	MongoUrl string
	DbName   string
	Db       *mgo.Database
}

func (api *GoodBadApi) Init() {
	session, err := mgo.Dial(api.MongoUrl)
	if err != nil {
		panic(err)
	}

	api.Db = session.DB(api.DbName)

	models.Init(api.Db)
}

func rateLimit(db *mgo.Database, imei, line string) error {
	dlp := models.FindDeviceLastPost(db, imei, line)

	neverPost := dlp.Timestamp == 0
	oldPost := (time.Now().Unix() - dlp.Timestamp) > int64((30 * time.Minute).Seconds())

	if neverPost || oldPost {
		// Register DeviceLastPost
		newDlp := &models.DeviceLastPost{
			Imei:      imei,
			Line:      line,
			Timestamp: time.Now().Unix(),
		}

		newDlp.Save(db)

		return nil
	}

	return errors.New("Rate limited")
}

func (api *GoodBadApi) PostGoodBad(w rest.ResponseWriter, req *rest.Request) {
	goodBad := models.GoodBad{}
	if err := req.DecodeJsonPayload(&goodBad); err != nil {
		panic(err)
	}

	if err := rateLimit(api.Db, goodBad.Imei, goodBad.Line); err != nil {
		w.WriteHeader(400)

		return
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

type Live struct {
	Status string
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
		&rest.Route{"GET", "/live", func(w rest.ResponseWriter, req *rest.Request) {
			w.WriteJson(&Live{
				Status: "I'm alive!",
			})
		}},
		rest.RouteObjectMethod("POST", "/good_bad", api, "PostGoodBad"),
	)
	http.ListenAndServe(":"+port, &handler)
}
