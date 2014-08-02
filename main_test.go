package main

import (
	"testing"

	"github.com/rodrigosaito/tabomtaruim/models"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/stretchr/testify/assert"
)

func prepareHandler() (rest.ResourceHandler, *GoodBadApi) {
	handler := rest.ResourceHandler{}

	api := &GoodBadApi{
		MongoUrl: "localhost",
		DbName:   "good_bad_test",
	}
	api.Init()
	api.Db.DropDatabase()

	handler.SetRoutes(
		rest.RouteObjectMethod("POST", "/good_bad", api, "PostGoodBad"),
	)

	return handler, api
}

func TestPostGoodBad(t *testing.T) {
	handler, _ := prepareHandler()

	recorded := test.RunRequest(t, &handler,
		test.MakeSimpleRequest("POST", "http://1.2.3.4/good_bad", models.GoodBad{"cptm-9", "123321", "good", 0}))

	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()

	lineStatus := &models.LineStatus{}

	recorded.DecodeJsonPayload(lineStatus)

	assert.Equal(t, "cptm-9", lineStatus.Line)
	assert.Equal(t, 10, lineStatus.Goods)
	assert.Equal(t, 2, lineStatus.Bads)
	assert.Equal(t, "good", lineStatus.Status)
}

func TestDoublePostGoodBad(t *testing.T) {
	handler, _ := prepareHandler()

	recorded1 := test.RunRequest(t, &handler,
		test.MakeSimpleRequest("POST", "http://1.2.3.4/good_bad", models.GoodBad{"cptm-9", "123321", "good", 0}))
	recorded1.CodeIs(200)

	recorded2 := test.RunRequest(t, &handler,
		test.MakeSimpleRequest("POST", "http://1.2.3.4/good_bad", models.GoodBad{"cptm-9", "123321", "bad", 0}))
	recorded2.CodeIs(400)
}

func TestPostGoodBadWithVeryOldDeviceLastPost(t *testing.T) {
	handler, api := prepareHandler()

	// Saves a really old DeviceLastPost to ensure that rate limit is working
	dlp := &models.DeviceLastPost{
		Imei:      "123321",
		Line:      "cptm-9",
		Timestamp: 123,
	}
	dlp.Save(api.Db)

	recorded := test.RunRequest(t, &handler,
		test.MakeSimpleRequest("POST", "http://1.2.3.4/good_bad", models.GoodBad{"cptm-9", "123321", "good", 0}))
	recorded.CodeIs(200)
}

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
