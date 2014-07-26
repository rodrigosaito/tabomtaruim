package models

import (
	"testing"

	// "github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	api.Db = session.DB("good_bad_test")
}
