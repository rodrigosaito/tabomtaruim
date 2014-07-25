package main

import (
	"testing"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/stretchr/testify/assert"
)

func prepareHandler() rest.ResourceHandler {
	handler := rest.ResourceHandler{}

	api := &GoodBadApi{
		MongoUrl: "localhost",
		DbName:   "good_bad_test",
	}
	api.Init()

	handler.SetRoutes(
		rest.RouteObjectMethod("POST", "/good_bad", api, "PostGoodBad"),
	)

	return handler
}

func TestPostGoodBad(t *testing.T) {
	handler := prepareHandler()

	recorded := test.RunRequest(t, &handler,
		test.MakeSimpleRequest("POST", "http://1.2.3.4/good_bad", GoodBad{"cptm-9", "123321", "", 0}))

	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()

	lineStatus := &LineStatus{}

	recorded.DecodeJsonPayload(lineStatus)

	assert.Equal(t, "cptm-9", lineStatus.Line)
	assert.Equal(t, 10, lineStatus.Goods)
	assert.Equal(t, 2, lineStatus.Bads)
	assert.Equal(t, "good", lineStatus.Status)
}
