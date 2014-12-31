package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/rodrigosaito/tabomtaruim/models"
)

type RecordGoodBadApi struct {
}

func (api *RecordGoodBadApi) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	goodBad := models.GoodBad{}
	if err := DecodeJsonPayload(req, &goodBad); err != nil {
		panic(err)
	}

	if err := goodBad.Save(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	lineStatus := models.GetLineStatus(goodBad.Line)

	w.WriteHeader(http.StatusCreated)
	WriteJson(w, lineStatus)
}

type CheckStatusApi struct {
}

func (api *CheckStatusApi) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	line := req.URL.Query().Get(":line")

	lineStatus := models.GetLineStatus(line)

	WriteJson(w, lineStatus)
}

func DecodeJsonPayload(req *http.Request, v interface{}) error {
	content, err := ioutil.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, v)
	if err != nil {
		return err
	}
	return nil
}

func WriteJson(w http.ResponseWriter, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = w.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func LiveCheck(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "I'm alive!!!")
}

func main() {
	// configuration
	var port = flag.String("port", "8080", "port to run server")
	var mongoUrl = flag.String("mongo-url", "localhost", "url to connect to mongo server")
	var mongoDbName = flag.String("mongo-db-name", "good_bad_dev", "name of database to persist MongoDB documents")
	flag.Parse()

	// mongodb connection
	models.Init(*mongoUrl, *mongoDbName)

	// http server
	m := pat.New()
	m.Get("/live", http.HandlerFunc(LiveCheck))
	m.Post("/good_bad", &RecordGoodBadApi{})
	m.Get("/good_bad/:line", &CheckStatusApi{})

	http.Handle("/", m)

	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
