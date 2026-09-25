package oidcdot

import (
	"github.com/google/wire"
	oidcrp "github.com/scryinfo/dot/line/oidcdot/oidc_rp"
	"github.com/scryinfo/dot/line/oidcdot/oidc_server/oidc_storage"
)

var OidcPebble2Set = wire.NewSet(
	oidc_storage.Pebble2Set,
	oidcrp.NewAppSessionHelper,
)
