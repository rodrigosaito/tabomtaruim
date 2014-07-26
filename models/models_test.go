package models

import (
	"testing"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	db := session.DB("good_bad_test")

	before := GoodBadCount(db)

	goodBad := GoodBad{
		Line:      "cptm-9",
		Imei:      "123321",
		Status:    "good",
		Timestamp: time.Now().Unix(),
	}

	goodBad.Save(db)

	after := GoodBadCount(db)
	assert.Equal(t, before+1, after)
}
