package models

import (
	"testing"

	"gopkg.in/mgo.v2"

	"github.com/stretchr/testify/assert"
)

func Db() *mgo.Database {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	db := session.DB("good_bad_test")
	db.DropDatabase()

	Init(db)

	return db
}

func TestGoodBadSave(t *testing.T) {
	db := Db()
	Init(db)

	before := GoodBadCount()

	goodBad := GoodBad{
		Line:   "cptm-9",
		Imei:   "123321",
		Status: "good",
	}

	goodBad.Save()

	after := GoodBadCount()
	assert.Equal(t, before+1, after)

	if goodBad.Timestamp == 0 {
		t.Error("Should set current timestamp before saving")
	}
}
