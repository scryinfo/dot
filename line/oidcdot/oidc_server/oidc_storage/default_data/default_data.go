package oidc_storage

import (
	_ "embed"
	"encoding/json/v2"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/line/db/badgerdot"
	"github.com/scryinfo/dot/line/db/pebble2dot"
	"github.com/scryinfo/dot/line/oidcdot/oidc_server/oidc_storage"
)

type DefaultDataPebble2 struct {
}

func NewDefaultDataPebble2(db *pebble2dot.Pebble2, logger *dot.LoggerType) (*DefaultDataPebble2, error) {
	_, err := InitPebble2(db, logger)
	return &DefaultDataPebble2{}, err
}

type DefaultDataBadger struct {
}

func NewDefaultDataBadger(db *badgerdot.BadgerDbDot, logger *dot.LoggerType) (*DefaultDataBadger, error) {
	_, err := InitBadger(db, logger)
	return &DefaultDataBadger{}, err
}

type DefaultData struct {
	Clients []oidc_storage.OidcClient `json:"clients"`
	Users   []oidc_storage.User       `json:"users"`
}

//go:embed default_data.json
var configData []byte
var defaultData DefaultData

func Get() (*DefaultData, error) {
	if err := json.Unmarshal(configData, &defaultData); err != nil {
		return nil, err
	}
	return &defaultData, nil
}

func InitPebble2(db *pebble2dot.Pebble2, logger *dot.LoggerType) (*DefaultData, error) {
	_, err := Get()
	if err != nil {
		return nil, err
	}
	userDao := oidc_storage.NewUserDaoPebble2(db, logger)
	for i, _ := range defaultData.Users {
		userDao.Add(&defaultData.Users[i])
	}
	clientDao := oidc_storage.NewOidcClientDaoPebble2(db, logger)
	for i, _ := range defaultData.Clients {
		clientDao.Add(&defaultData.Clients[i])
	}
	return &defaultData, nil
}

func InitBadger(db *badgerdot.BadgerDbDot, logger *dot.LoggerType) (*DefaultData, error) {
	_, err := Get()
	if err != nil {
		return nil, err
	}
	userDao := oidc_storage.NewUserDaoBadger(db, logger)
	for i, _ := range defaultData.Users {
		userDao.Add(&defaultData.Users[i])
	}
	clientDao := oidc_storage.NewOidcClientDaoBadger(db, logger)
	for i, _ := range defaultData.Clients {
		clientDao.Add(&defaultData.Clients[i])
	}
	return &defaultData, nil
}
