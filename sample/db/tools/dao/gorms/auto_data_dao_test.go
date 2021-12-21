package gorms

import (
	"github.com/scryinfo/dot/dots/db/gorms"
	"github.com/scryinfo/dot/sample/db/tools/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

const dsn = "root:123456@tcp(127.0.0.1:3306)/dot_test?charset=utf8mb4&parseTime=True&loc=Local"

func CreateTestConn() *AutoDataDao {
	db := gorms.NewGormsTest(dsn, true, "mysql")
	if db != nil {
		base := gorms.GenerateDaoBaseByDb(db)
		return &AutoDataDao{base}
	}
	return nil
}

func TestCreateTable(t *testing.T) {
	dao := CreateTestConn()
	dao.Wrapper.GetDb().AutoMigrate(&model.AutoData{})
}

func TestAutoDataDao_GetByID(t *testing.T) {
	dao := CreateTestConn()
	m, err := dao.GetByID(1)
	if err != nil {
		return
	}
	assert.Equal(t, m.Age, 12)
}

func TestAutoDataDao_Query(t *testing.T) {
	dao := CreateTestConn()
	ms, err := dao.Query("age <  ? ", 25)
	if err != nil {
		return
	}
	assert.Equal(t, len(ms), 0)
}
func TestAutoDataDao_List(t *testing.T) {
	dao := CreateTestConn()
	ms, err := dao.List()
	if err != nil {
		return
	}
	assert.Equal(t, len(ms), 5)
}
func TestAutoDataDao_Count(t *testing.T) {
	dao := CreateTestConn()
	count, err := dao.Count("age > ? ", int8(24))
	if err != nil {
		return
	}
	assert.Equal(t, count, int64(4))
}
func TestAutoDataDao_QueryPage(t *testing.T) {
	dao := CreateTestConn()
	ms, err := dao.QueryPage(2, 2, "age > ?", 24)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(ms), 2)
}
func TestAutoDataDao_QueryPageWithCount(t *testing.T) {
	dao := CreateTestConn()
	ms, count, err := dao.QueryPageWithCount(2, 2, "age > ?", 24)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(ms), 2)
	assert.Equal(t, count, int64(2))

}
func TestAutoDataDao_QueryOne(t *testing.T) {
	dao := CreateTestConn()
	m, err := dao.QueryOne("age > ? and age <= ? ", 12, 25)
	assert.Equal(t, err, nil)
	assert.Equal(t, m.ID, uint(1))
}

func TestAutoDataDao_Insert(t *testing.T) {
	dao := CreateTestConn()
	m := &model.AutoData{
		Name: "jayce",
		Age:  25,
	}
	assert.Equal(t, dao.Insert(m), nil)
}

func TestAutoDataDao_Inserts(t *testing.T) {
	dao := CreateTestConn()
	ms := make([]*model.AutoData, 0)
	ms = append(ms, &model.AutoData{
		Name: "jayce",
		Age:  20,
	})
	ms = append(ms, &model.AutoData{
		Name: "jayce",
		Age:  21,
	})
	ms = append(ms, &model.AutoData{
		Name: "jayce",
		Age:  22,
	})
	assert.Equal(t, dao.Inserts(ms), nil)
}

func TestAutoDataDao_Upsert(t *testing.T) {
	dao := CreateTestConn()
	m := &model.AutoData{
		Name: "herry",
	}
	//m.ID = 1
	assert.Equal(t, dao.Upsert(m), nil)

}

func TestAutoDataDao_Upserts(t *testing.T) {
	dao := CreateTestConn()
	ms := make([]*model.AutoData, 0)
	ms = append(ms, &model.AutoData{
		Name: "jayce",
		Age:  20,
	})
	ms = append(ms, &model.AutoData{
		Name: "jayce",
		Age:  21,
	})
	ms = append(ms, &model.AutoData{
		Name: "jayce",
		Age:  22,
	})
	ms[0].ID = 1
	ms[1].ID = 11
	assert.Equal(t, dao.Upserts(ms), nil)

}

func TestAutoDataDao_Update(t *testing.T) {
	dao := CreateTestConn()
	m := &model.AutoData{
		Name: "Herry",
	}
	m.ID = 1
	assert.Equal(t, dao.Update(m), nil)

}

func TestAutoDataDao_Save(t *testing.T) {
	dao := CreateTestConn()
	m := &model.AutoData{
		Name: "Herry",
	}
	m.ID = 1
	assert.Equal(t, dao.Save(m), nil)
	assert.Equal(t, m.Age, int8(0))

}

func TestAutoDataDao_UpdateColumn(t *testing.T) {
	dao := CreateTestConn()
	m := &model.AutoData{}
	m.ID = 1
	assert.Equal(t, dao.UpdateColumn(m, "age", 12), nil)
}

func TestAutoDataDao_DeleteById(t *testing.T) {
	dao := CreateTestConn()
	assert.Equal(t, dao.DeleteById(1), nil)
}

func TestAutoDataDao_DeleteByIds(t *testing.T) {
	dao := CreateTestConn()
	assert.Equal(t, dao.DeleteByIds([]uint{2, 3}), nil)

}

func TestAutoDataDao_Delete(t *testing.T) {
	dao := CreateTestConn()
	assert.Equal(t, dao.Delete("age < ?", 25), nil)
}

func TestAutoDataDao_DeleteByIdUnscoped(t *testing.T) {
	dao := CreateTestConn()
	assert.Equal(t, dao.DeleteByIdUnscoped(1), nil)
}
