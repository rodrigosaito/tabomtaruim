package main

import (
	"testing"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"
)

func TestPostGoodBad(t *testing.T) {
	handler := rest.ResourceHandler{}
	handler.SetRoutes(
		&rest.Route{"POST", "/good_bad", PostGoodBad},
	)

	recorded := test.RunRequest(t, &handler,
		test.MakeSimpleRequest("POST", "http://1.2.3.4/good_bad", nil))

	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
}
