package pgs

type ModelBase struct {
	ID                    string `pg:",pk" json:"id"`                          //id
	UpdateTime            int64  `json:"updateTime"`                           //更新时间
	CreateTime            int64  `json:"createTime"`                           //创建时间
	OptimisticLockVersion int64  `pg:",use_zero" json:"optimisticLockVersion"` //default 0 ，not null
}
