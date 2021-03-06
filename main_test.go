package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rodrigosaito/tabomtaruim/models"
	"github.com/sfreiberg/mongo"
	"github.com/stretchr/testify/assert"
)

func init() {
	dbName := "good_bad_test_main"

	models.Init("localhost", dbName)

	session, _ := mongo.GetSession()
	session.DB(dbName).DropDatabase()
}

func postGoodBadTest(v interface{}) (*models.LineStatus, *httptest.ResponseRecorder) {
	b, _ := json.Marshal(v)
	req, err := http.NewRequest("POST", "http://example.com/good_bad", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()

	api := RecordGoodBadApi{}
	api.ServeHTTP(w, req)

	lineStatus := &models.LineStatus{}

	json.Unmarshal(w.Body.Bytes(), lineStatus)

	return lineStatus, w
}

func TestPostGoodBad(t *testing.T) {
	lineStatus, w := postGoodBadTest(&models.GoodBad{"cptm-9", "123", "good", 0})

	assert.Equal(t, http.StatusCreated, w.Code)

	assert.Equal(t, "cptm-9", lineStatus.Line)
	assert.Equal(t, 1, lineStatus.Goods)
	assert.Equal(t, 0, lineStatus.Bads)
	assert.Equal(t, "good", lineStatus.Status)
}

func TestDoublePostGoodBad(t *testing.T) {
	t.SkipNow() // Skipping this test until rate limit is implemented

	goodBad := &models.GoodBad{"cptm-9", "1234567890", "good", 0}

	_, w := postGoodBadTest(goodBad)
	assert.Equal(t, 201, w.Code) // assert that the first request succeeded

	_, w = postGoodBadTest(goodBad)

	assert.Equal(t, 403, w.Code)
}

/*
func TestDoublePostGoodBadWithVeryOldDeviceLastPost(t *testing.T) {
	handler, api := prepareHandler()

	// Saves a really old DeviceLastPost to ensure that rate limit is working
	dlp := &models.DeviceLastPost{
		Imei:      "123321",
		Line:      "cptm-9",
		Timestamp: 123,
	}
	dlp.Save(api.Db)

	recorded1 := test.RunRequest(t, &handler,
		test.MakeSimpleRequest("POST", "http://1.2.3.4/good_bad", models.GoodBad{"cptm-9", "123321", "good", 0}))
	recorded1.CodeIs(200)

	recorded2 := test.RunRequest(t, &handler,
		test.MakeSimpleRequest("POST", "http://1.2.3.4/good_bad", models.GoodBad{"cptm-9", "123321", "bad", 0}))
	recorded2.CodeIs(400)
}
*/
