package domyApi

import (
	"os"

	"github.com/domyapi/helper/atdb"
)

var MongoString string = os.Getenv("MONGOSTRING")

var mongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "domyid",
}

var Mongoconn, ErrorMongoconn = atdb.MongoConnect(mongoinfo)
