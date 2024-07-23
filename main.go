package domyApi

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	route "github.com/domyid/domyapi/route"
)

func init() {
	functions.HTTP("WebHook", route.URL)
}
