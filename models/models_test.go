package models

import (
	"testing"

	"gopkg.in/mgo.v2"

	// "github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	db := session.DB("good_bad_test")

	db.C("good_bad")
}
