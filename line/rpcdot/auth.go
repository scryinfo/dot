package rpcdot

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/scryinfo/dot/dot"
	httptools "github.com/scryinfo/dot/line/rpcdot/http_tools"
)

func NewHandlerMiddle() HandlerMiddle {
	return authWare
}

var unauthUrls = []string{}

func authWare(_ http.ResponseWriter, r *http.Request) error {
	if slices.Contains(unauthUrls, r.URL.Path) {
		return nil
	} else if len(unauthUrls) == 1 && unauthUrls[0] == "*" {
		return nil
	}
	tokenString := httptools.GetToken(r)
	if len(tokenString) < 1 {
		err := fmt.Errorf("Unauthorized")
		dot.Logger.Info().Err(err).Send()
		return err
	}
	return httptools.ValidToken(tokenString)
}
