package oidc_storage

import (
	"context"

	daobase "github.com/scryinfo/dot/line/db/dao/dao_base"
)

func (s *StoragePebble2) UserByID(ctx context.Context, userId string) (User, error) {
	return s.userDao.Find(daobase.IdType(userId))
}
func (s *StoragePebble2) UserByKeyname(ctx context.Context, keyname string) (*User, error) {
	return s.userDao.FindByKeyname(keyname)
}
