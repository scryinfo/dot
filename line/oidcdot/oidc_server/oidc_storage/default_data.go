package oidc_storage

import (
	_ "embed"
	"encoding/json"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/db/badgerdot"
	"github.com/scryinfo/dot/line/db/pebble2dot"
)

type DefaultDataPebble2 struct {
}

func NewDefaultDataPebble2(db *pebble2dot.Pebble2, logger *dot.LoggerType) DefaultDataPebble2 {
	{
		userDao := NewUserDaoPebble2(db, logger)
		for i, _ := range defaultData.Users {
			userDao.Add(&defaultData.Users[i])
		}
		clientDao := NewOidcClientDaoPebble2(db, logger)
		for i, _ := range defaultData.Clients {
			clientDao.Add(&defaultData.Clients[i])
		}
	}

	return DefaultDataPebble2{}
}

type DefaultDataBadger struct {
}

func NewDefaultDataBadger(db *badgerdot.BadgerDbDot, logger *dot.LoggerType) DefaultDataBadger {
	{
		userDao := NewUserDaoBadger(db, logger)
		for i, _ := range defaultData.Users {
			userDao.Add(&defaultData.Users[i])
		}
		clientDao := NewOidcClientDaoBadger(db, logger)
		for i, _ := range defaultData.Clients {
			clientDao.Add(&defaultData.Clients[i])
		}
	}
	return DefaultDataBadger{}
}

type DefaultData struct {
	Clients []OidcClient `json:"clients"`
	Users   []User       `json:"users"`
}

//go:embed default_data.json
var configData []byte
var defaultData DefaultData

func init() {
	if err := json.Unmarshal(configData, &defaultData); err != nil {
		panic(err)
	}
}
