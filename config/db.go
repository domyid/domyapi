package domyApi

import (
	"os"

	atdb "github.com/domyid/domyapi/helper/atdb"
)

var MongoString string = os.Getenv("MONGOSTRING")

var mongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "siakad",
}

var Mongoconn, ErrorMongoconn = atdb.MongoConnect(mongoinfo)
