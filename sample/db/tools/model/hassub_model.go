package model

import (
	"fmt"
	"github.com/scryinfo/dot/dots/db/pgs"
)

const (
	HasSub_Table                 = "has_subs"
	HasSub_Struct                = "has_sub"
	HasSub_ID                    = "id"
	HasSub_UpdateTime            = "update_time"
	HasSub_CreateTime            = "create_time"
	HasSub_OptimisticLockVersion = "optimistic_lock_version"
	HasSub_SubId                 = "sub_id"
	HasSub_Count                 = "count"
	HasSub_SubData               = "sub_data"
)

func (m *HasSub) String() string {
	//todo please change the format string
	//m.ID, m.UpdateTime, m.CreateTime, m.OptimisticLockVersion, m.SubId, m.Count, m.SubData,
	str := fmt.Sprintf("HasSub<%s %s >",
		m.ID, m.SubId,
	)
	return str
}

func (m *HasSub) ToMap() map[string]string {
	res := pgs.ToMap(m, map[string]bool{})
	return res
}

func (m *HasSub) ToUpsertSet() []string {
	res := []string{
		fmt.Sprintf("%s = EXCLUDED.%s", HasSub_ID, HasSub_ID),
		fmt.Sprintf("%s = EXCLUDED.%s", HasSub_UpdateTime, HasSub_UpdateTime),
		fmt.Sprintf("%s = EXCLUDED.%s", HasSub_CreateTime, HasSub_CreateTime),
		fmt.Sprintf("%s = EXCLUDED.%s", HasSub_OptimisticLockVersion, HasSub_OptimisticLockVersion),
		fmt.Sprintf("%s = EXCLUDED.%s", HasSub_SubId, HasSub_SubId),
		fmt.Sprintf("%s = EXCLUDED.%s", HasSub_Count, HasSub_Count),
	}
	return res
}
