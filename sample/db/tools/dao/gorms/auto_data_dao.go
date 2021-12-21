package gorms

import (
	"github.com/scryinfo/dot/dot"
	"github.com/scryinfo/dot/dots/db/gorms"
	"github.com/scryinfo/dot/dots/db/pgs"
	"github.com/scryinfo/dot/sample/db/tools/model"
	"gorm.io/gorm/clause"
)

const AutoDataDaoTypeID = "b7ea29d4-64d9-4ac1-a288-ecfcfbaa9fc6"

type AutoDataDao struct {
	*gorms.DaoBase `dot:""`
}

//AutoDataDaoTypeLives
func AutoDataDaoTypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeID: AutoDataDaoTypeID, NewDoter: func(conf []byte) (dot.Dot, error) {
			return &AutoDataDao{}, nil
		}},
		Lives: []dot.Live{
			{
				LiveID: AutoDataDaoTypeID,
				RelyLives: map[string]dot.LiveID{
					"DaoBase": pgs.DaoBaseTypeID,
				},
			},
		},
	}

	lives := gorms.DaoBaseTypeLives()
	lives = append(lives, tl)

	return lives
}

func (c *AutoDataDao) GetByID(id uint) (m *model.AutoData, err error) {
	result := c.Wrapper.GetDb().First(&m, id)
	if result.RowsAffected == 1 {
		return m, nil
	}
	return nil, result.Error
}

func (c *AutoDataDao) Query(condition string, params ...interface{}) (ms []*model.AutoData, err error) {

	db := c.Wrapper.GetDb()
	if len(condition) < 1 {
		db = db.Find(&ms)
	} else {
		db = db.Where(condition, params...).Find(&ms)
	}
	err = db.Error
	if err != nil {
		ms = nil
	}
	return
}

func (c *AutoDataDao) List() (ms []*model.AutoData, err error) {
	db := c.Wrapper.GetDb().Find(&ms)
	if db.Error != nil {
		ms = nil
	}
	return ms, db.Error
}

func (c *AutoDataDao) Count(condition string, params ...interface{}) (count int64, err error) {

	var ms []*model.AutoData
	db := c.Wrapper.GetDb()
	if len(condition) < 1 {
		db = db.Find(&ms)
	} else {
		db = db.Where(condition, params...).Find(&ms)
	}
	err = db.Error
	if err != nil {
		count = 0
	}
	count = db.RowsAffected
	return
}

func (c *AutoDataDao) QueryPage(pageSize int, page int, condition string, params ...interface{}) (ms []*model.AutoData, err error) {
	db := c.Wrapper.GetDb()
	if len(condition) < 1 {
		db = db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&ms)
	} else {
		db = db.Limit(pageSize).Offset((page-1)*pageSize).Where(condition, params...).Find(&ms)
	}
	if db.Error != nil { //be sure
		ms = nil
	}
	return ms, db.Error
}

func (c *AutoDataDao) QueryPageWithCount(pageSize int, page int, condition string, params ...interface{}) (ms []*model.AutoData, count int64, err error) {
	db := c.Wrapper.GetDb()
	if len(condition) < 1 {
		db = db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&ms)
	} else {
		db = db.Limit(pageSize).Offset((page-1)*pageSize).Where(condition, params...).Find(&ms)
	}
	if db.Error != nil { //be sure
		ms = nil
	}
	return ms, db.RowsAffected, db.Error
}

func (c *AutoDataDao) QueryOne(condition string, params ...interface{}) (m *model.AutoData, err error) {
	db := c.Wrapper.GetDb()
	if len(condition) < 1 {
		db = db.First(&m)
	} else {
		db = db.Where(condition, params...).First(&m)
	}
	err = db.Error
	if err != nil {
		m = nil
	}
	return
}

//insert = insertReturn
func (c *AutoDataDao) Insert(m *model.AutoData) (err error) {
	//because id auto increment
	//m.ID = 0
	err = c.Wrapper.GetDb().Create(&m).Error
	return
}

func (c *AutoDataDao) Inserts(ms []*model.AutoData) (err error) {
	//because id auto increment
	//m.ID = 0
	err = c.Wrapper.GetDb().Create(&ms).Error
	return
}

//update everything except ID
func (c *AutoDataDao) Upsert(m *model.AutoData) (err error) {
	//because id auto increment
	//m.ID = 0
	err = c.Wrapper.GetDb().Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&m).Error
	return
}

func (c *AutoDataDao) Upserts(ms []*model.AutoData) (err error) {
	//because id auto increment
	//m.ID = 0
	err = c.Wrapper.GetDb().Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&ms).Error
	return
}

//Updates 方法支持 struct 和 map[string]interface{} 参数
//当使用 struct 更新时，默认情况下，GORM 只会更新非零值的字段
func (c *AutoDataDao) Update(m *model.AutoData) (err error) {
	err = c.Wrapper.GetDb().Updates(&m).Error
	return
}

//Save 会保存所有的字段，即使字段是零值
//todo CreateTime=0
func (c *AutoDataDao) Save(m *model.AutoData) (err error) {
	err = c.Wrapper.GetDb().Save(&m).Error
	return
}

//default where id=m.ID
func (c *AutoDataDao) UpdateColumn(m *model.AutoData, columnName string, value interface{}) (err error) {
	err = c.Wrapper.GetDb().Model(&m).Update(columnName, value).Error
	return
}

//soft delete
func (c *AutoDataDao) DeleteById(id uint) error {
	return c.Wrapper.GetDb().Delete(&model.AutoData{}, id).Error
}

//soft delete
func (c *AutoDataDao) DeleteByIds(ids []uint) error {
	return c.Wrapper.GetDb().Delete(&model.AutoData{}, ids).Error
}

//soft delete
func (c *AutoDataDao) Delete(condition string, params ...interface{}) error {
	return c.Wrapper.GetDb().Where(condition, params...).Delete(&model.AutoData{}).Error
}

//Delete permanently
func (c *AutoDataDao) DeleteByIdUnscoped(id uint) error {
	return c.Wrapper.GetDb().Unscoped().Delete(&model.AutoData{}, id).Error
}
