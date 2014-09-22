package models

import (
	"testing"

	"gopkg.in/mgo.v2"

	"github.com/stretchr/testify/assert"

	// "gopkg.in/mgo.v2/bson"
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

func TestGoodCountBefore30Minutes(t *testing.T) {
	t.SkipNow()
	db := Db()
	Init(db)

	before := Count("good", "cptm-9")

	goodBad := GoodBad{
		Line:   "cptm-9",
		Imei:   "123321",
		Status: "good",
	}

	goodBad.Save()

	// db.C("good_bad").UpdateAll(bson.M{}, bson.M{"timestamp": "1411424203"})

	after := GoodBadCount()
	assert.Equal(t, before, after)

	if goodBad.Timestamp == 0 {
		t.Error("Should set current timestamp before saving")
	}
}
