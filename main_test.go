package main

import (
	"testing"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/stretchr/testify/assert"
)

func prepareHandler() rest.ResourceHandler {
	handler := rest.ResourceHandler{}
	handler.SetRoutes(
		&rest.Route{"POST", "/good_bad", PostGoodBad},
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
