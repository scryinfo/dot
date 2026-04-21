package db

type ModelBase struct {
	ID                    string `bun:"id,pk" json:"id"`        //id
	UpdateTime            int64  `json:"updateTime"`            //更新时间
	CreateTime            int64  `json:"createTime"`            //创建时间
	OptimisticLockVersion int64  `json:"optimisticLockVersion"` //default 0 ，not null
}
