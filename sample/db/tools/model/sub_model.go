package model

import (
	"fmt"
	"github.com/scryinfo/dot/dots/db/pgs"
)

const (
	Sub_Table                 = "subs"
	Sub_Struct                = "sub"
	Sub_ID                    = "id"
	Sub_UpdateTime            = "update_time"
	Sub_CreateTime            = "create_time"
	Sub_OptimisticLockVersion = "optimistic_lock_version"
	Sub_Name                  = "name"
)

func (m *Sub) String() string {
	//todo please change the format string
	//m.ID, m.UpdateTime, m.CreateTime, m.OptimisticLockVersion, m.Name,
	str := fmt.Sprintf("Sub<%s %s >",
		m.ID, m.Name,
	)
	return str
}

func (m *Sub) ToMap() map[string]string {
	res := pgs.ToMap(m, map[string]bool{})
	return res
}

func (m *Sub) ToUpsertSet() []string {
	res := []string{
		fmt.Sprintf("%s = EXCLUDED.%s", Sub_ID, Sub_ID),
		fmt.Sprintf("%s = EXCLUDED.%s", Sub_UpdateTime, Sub_UpdateTime),
		fmt.Sprintf("%s = EXCLUDED.%s", Sub_CreateTime, Sub_CreateTime),
		fmt.Sprintf("%s = EXCLUDED.%s", Sub_OptimisticLockVersion, Sub_OptimisticLockVersion),
		fmt.Sprintf("%s = EXCLUDED.%s", Sub_Name, Sub_Name),
	}
	return res
}
