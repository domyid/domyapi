package domyApi

import (
	"log"
	"strings"

	sy "github.com/JPratama7/util/sync"
	at "github.com/domyid/domyapi/helper/at"
	atdb "github.com/domyid/domyapi/helper/atdb"
	model "github.com/domyid/domyapi/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	PoolStringBuilder *sy.Pool[*strings.Builder]
)

var WAAPIToken string

var IPPort, Net = at.GetAddress()

var PhoneNumber string

func SetEnv() {
	if ErrorMongoconn != nil {
		log.Println(ErrorMongoconn.Error())
	}
	profile, err := atdb.GetOneDoc[model.Profile](Mongoconn, "profile", primitive.M{})
	if err != nil {
		log.Println(err)
	}

	WAAPIToken = profile.Token
}
