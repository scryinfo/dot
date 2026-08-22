package oidc_storage

import (
	"encoding/json/v2"
	"fmt"
	"testing"

	"github.com/scryinfo/dot/lib/kits"
	"github.com/scryinfo/dot/line/oidcdot/oidc_server/oidc_storage"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestDefaultData(t *testing.T) {
	data := DefaultData{
		Clients: make([]oidc_storage.OidcClient, 0, 2),
		Users:   make([]oidc_storage.User, 0, 4),
	}
	{
		data.Clients = append(data.Clients, oidc_storage.OidcClient{
			Id:           kits.Ids.NewXId(),
			Secret:       "66",
			ShowName:     "client1",
			RedirectUris: []string{"http://localhost:8089/callback"},
		})
		data.Clients = append(data.Clients, oidc_storage.OidcClient{
			Id:           kits.Ids.NewXId(),
			Secret:       "66",
			ShowName:     "client2",
			RedirectUris: []string{"http://localhost:8089/callback"},
		})
		data.Users = append(data.Users, oidc_storage.User{
			Id:                kits.Ids.NewXId(),
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
			Id:                kits.Ids.NewXId(),
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
			Id:                kits.Ids.NewXId(),
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
			Id:                kits.Ids.NewXId(),
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
			p, err := oidc_storage.PasswordHash.HashPassword(c.Secret)
			assert.Nil(t, err)
			c.Secret = p
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
}
