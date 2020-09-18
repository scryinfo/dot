package model

import (
	"fmt"
	"github.com/scryinfo/dot/dots/db/pgs"
)

const (
	Notice_Table                 = "notice"
	Notice_Struct                = "notice"
	Notice_ID                    = "id"
	Notice_Status                = "status"
	Notice_CreateTime            = "create_time"
	Notice_UpdateTime            = "update_time"
	Notice_OptimisticLockVersion = "optimistic_lock_version"
)

func (m *Notice) String() string {
	//todo please change the format string
	//m.ID, m.Status, m.CreateTime, m.UpdateTime, m.OptimisticLockVersion,
	str := fmt.Sprintf("Notice<%s >",
		m.ID,
	)
	return str
}

func (m *Notice) ToMap() map[string]string {
	res := pgs.ToMap(m, map[string]bool{})
	return res
}

func (m *Notice) ToUpsertSet() []string {
	res := []string{

		fmt.Sprintf("%s = EXCLUDED.%s", Notice_ID, Notice_ID),
		fmt.Sprintf("%s = EXCLUDED.%s", Notice_Status, Notice_Status),
		fmt.Sprintf("%s = EXCLUDED.%s", Notice_CreateTime, Notice_CreateTime),
		fmt.Sprintf("%s = EXCLUDED.%s", Notice_UpdateTime, Notice_UpdateTime),
		fmt.Sprintf("%s = EXCLUDED.%s", Notice_OptimisticLockVersion, Notice_OptimisticLockVersion),
	}
	return res
}
