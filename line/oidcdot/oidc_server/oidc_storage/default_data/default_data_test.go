package default_data

import (
	"encoding/json/v2"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/lib/kits"
	"github.com/scryinfo/dot/line/db/badgerdot"
	daobase "github.com/scryinfo/dot/line/db/dao/dao_base"
	"github.com/scryinfo/dot/line/db/pebble2dot"
	oidcapiv1 "github.com/scryinfo/dot/line/oidcdot/oidc_gen/oidcapi/v1"
	"github.com/scryinfo/dot/line/oidcdot/oidc_server/oidc_storage"
	"github.com/scryinfo/dot/line/sconfig"
	"github.com/scryinfo/scryg/sutils/sfile"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestGenerateDefaultData(t *testing.T) {
	data := DefaultData{
		Clients: make([]oidc_storage.OidcClient, 0, 2),
		Users:   make([]oidc_storage.User, 0, 4),
	}
	{
		data.Clients = append(data.Clients, oidc_storage.OidcClient{
			Id:               "dajma87ip8rh292hmrk0",
			SecretF:          "66",
			ShowNameF:        "client1",
			ApplicationTypeF: oidcapiv1.ApplicationType_APPLICATION_TYPE_WEB,
			RedirectUrisF:    []string{"http://localhost:8089/auth2/callback"},
			ResponseTypesF:   oidc_storage.ResponseTypeEx.EnumsProto(),
			GrantTypesF:      oidc_storage.GrantTypeEx.EnumsProto(),
		})
		data.Clients = append(data.Clients, oidc_storage.OidcClient{
			Id:               "dajma87ip8rh292hmrkg",
			SecretF:          "66",
			ShowNameF:        "client2",
			ApplicationTypeF: oidcapiv1.ApplicationType_APPLICATION_TYPE_WEB,
			RedirectUrisF:    []string{"http://localhost:8089/auth2/callback"},
			ResponseTypesF:   oidc_storage.ResponseTypeEx.EnumsProto(),
			GrantTypesF:      oidc_storage.GrantTypeEx.EnumsProto(),
		})
		data.Users = append(data.Users, oidc_storage.User{
			Id:                "dajma87ip8rh292hmrl0",
			Username:          "test1",
			Password:          "66",
			FirstName:         "test1-1",
			LastName:          "test1-2",
			Email:             "test1@scryinfo.info",
			EmailVerified:     false,
			Phone:             "test1-phone",
			PhoneVerified:     false,
			PreferredLanguage: language.English.String(),
		})
		data.Users = append(data.Users, oidc_storage.User{
			Id:                "dajma87ip8rh292hmrlg",
			Username:          "test2",
			Password:          "66",
			FirstName:         "test2-1",
			LastName:          "test2-2",
			Email:             "test2@scryinfo.info",
			EmailVerified:     false,
			Phone:             "test2-phone",
			PhoneVerified:     false,
			PreferredLanguage: language.English.String(),
		})
		data.Users = append(data.Users, oidc_storage.User{
			Id:                "dajma87ip8rh292hmrm0",
			Username:          "test3",
			Password:          "66",
			FirstName:         "test3-1",
			LastName:          "test3-2",
			Email:             "test3@scryinfo.info",
			EmailVerified:     false,
			Phone:             "test3-phone",
			PhoneVerified:     false,
			PreferredLanguage: language.English.String(),
		})
		data.Users = append(data.Users, oidc_storage.User{
			Id:                "dajma87ip8rh292hmrmg",
			Username:          "test4",
			Password:          "66",
			FirstName:         "test4-1",
			LastName:          "test4-2",
			Email:             "test4@scryinfo.info",
			EmailVerified:     false,
			Phone:             "test4-phone",
			PhoneVerified:     false,
			PreferredLanguage: language.English.String(),
		})
		for i, _ := range data.Clients {
			c := &data.Clients[i]
			p, err := oidc_storage.PasswordHash.HashPassword(c.SecretF)
			assert.Nil(t, err)
			c.SecretF = p
		}
		for i, _ := range data.Users {
			u := &data.Users[i]
			p, err := oidc_storage.PasswordHash.HashPassword(u.Password)
			assert.Nil(t, err)
			u.Password = p
		}
	}
	bs, err := json.Marshal(data)
	assert.Nil(t, err)
	fmt.Println(string(bs))
	fmt.Printf("\n\n")
}

func TestInitPebble2(t *testing.T) {
	db, fClear, err := newTestPebble2(t)
	defer fClear()
	assert.Nil(t, err)
	InitPebble2(db, logger)
	defaultData, err := Get()
	assert.Nil(t, err)
	{
		userDao := oidc_storage.NewUserDaoPebble2(db, logger)
		users, err := userDao.NextPage(&daobase.PageMeta{PageSize: int32(len(defaultData.Users) + 1)})
		assert.Nil(t, err)
		assert.Equal(t, len(defaultData.Users), len(users))
		assert.Equal(t, defaultData.Users, users)
	}
	{
		clientDao := oidc_storage.NewOidcClientDaoPebble2(db, logger)
		clients, err := clientDao.NextPage(&daobase.PageMeta{PageSize: int32(len(defaultData.Clients) + 1)})
		assert.Nil(t, err)
		assert.Equal(t, len(defaultData.Clients), len(clients))
		assert.Equal(t, defaultData.Clients, clients)
	}

}
func TestInitBadger(t *testing.T) {
	db, fClear, err := newTestBadger(t)
	defer fClear()
	assert.Nil(t, err)
	InitBadger(db, logger)
	defaultData, err := Get()
	assert.Nil(t, err)
	{
		userDao := oidc_storage.NewUserDaoBadger(db, logger)
		users, err := userDao.NextPage(&daobase.PageMeta{PageSize: int32(len(defaultData.Users) + 1)})
		assert.Nil(t, err)
		assert.Equal(t, len(defaultData.Users), len(users))
		assert.Equal(t, defaultData.Users, users)
	}
	{
		clientDao := oidc_storage.NewOidcClientDaoBadger(db, logger)
		clients, err := clientDao.NextPage(&daobase.PageMeta{PageSize: int32(len(defaultData.Clients) + 1)})
		assert.Nil(t, err)
		assert.Equal(t, len(defaultData.Clients), len(clients))
		assert.Equal(t, defaultData.Clients, clients)
	}
}

var logger = dot.NewTestLogger()

func newTestPebble2(t *testing.T) (*pebble2dot.Pebble2, func(), error) {
	sourcePath := filepath.Dir(kits.Config.GetCallSourceFile())
	config := pebble2dot.Pebble2Config{
		DbPath: filepath.Join(sourcePath, "temp/pebble"),
	}
	if !sfile.ExistDir(config.DbPath) {
		err := os.MkdirAll(config.DbPath, 0755)
		assert.Nil(t, err)
	}
	return pebble2dot.NewPebble2(&config, sconfig.NewTestSConfig(sourcePath, sourcePath, sourcePath), logger)
}
func newTestBadger(t *testing.T) (*badgerdot.BadgerDbDot, func(), error) {
	sourcePath := filepath.Dir(kits.Config.GetCallSourceFile())
	config := badgerdot.BadgerDbDotConfig{
		DbPath: filepath.Join(sourcePath, "temp/badger"),
	}
	if !sfile.ExistDir(config.DbPath) {
		err := os.MkdirAll(config.DbPath, 0755)
		assert.Nil(t, err)
	}
	return badgerdot.NewBadgerDot(&config, sconfig.NewTestSConfig(sourcePath, sourcePath, sourcePath), logger)
}
