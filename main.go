package domyApi

import (
	route "github.com/domyid/domyapi/route"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("WebHook", route.URL)
}
