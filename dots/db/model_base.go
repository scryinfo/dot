package db

import "gorm.io/gorm"

type ModelBase struct {
	ID                    string `bun:"id,pk" pg:",pk" json:"id"`                //id
	UpdateTime            int64  `json:"update_time"`                            //更新时间
	CreateTime            int64  `json:"create_time"`                            //创建时间
	OptimisticLockVersion int64  `pg:",use_zero" json:"optimistic_lock_version"` //default 0 ，not null
}

type AutoModelBase struct {
	ID          uint           `gorm:"primaryKey;autoIncrement"`
	CreatedTime int64          `gorm:"autoCreateTime"`
	UpdatedTime int64          `gorm:"autoCreateTime;autoUpdateTime"`
	DeletedTime gorm.DeletedAt `gorm:"index"`
}
