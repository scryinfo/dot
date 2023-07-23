package pg

func GetDaoData() string {
	return pgDaoData
}

const pgDaoData = `
package {{$.DaoPkgName}}
import (
	"context"
	"database/sql"
	"time"
	
    "github.com/scryinfo/dot/dot"
    "github.com/scryinfo/dot/dots/db/bun/pgd"
	"github.com/scryinfo/scryg/sutils/uuid"
	"github.com/uptrace/bun"
	"{{$.ImportModelPkgName}}"
)

const {{$.DaoName}}TypeID = "{{$.ID}}"

type {{$.DaoName}} struct {
	*pgd.DaoBase {{$.BackQuote}}dot:""{{$.BackQuote}}
}

//{{$.DaoName}}TypeLives
func {{$.DaoName}}TypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{
			Name: "{{$.DaoName}}",
			TypeID: {{$.DaoName}}TypeID, 
			NewDoter: func(conf []byte) (dot.Dot, error) {
				return &{{$.DaoName}}{}, nil
			},
		},
		Lives: []dot.Live{
			{
				LiveID: {{$.DaoName}}TypeID,
				RelyLives: map[string]dot.LiveID{
					"DaoBase": pgd.DaoBaseTypeID,
				},
			},
		},
	}

	lives := pgd.DaoBaseTypeLives()
	lives = append(lives, tl)
	return lives
}

// if find|update|insert nothing, sql.ErrNoRows error may returned

func (c *{{$.DaoName}}) GetByIDWithLock(conn bun.IDB, id string) (m *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	m = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	m.ID = id
	err = conn.NewSelect().Model(m).WherePK().For("UPDATE").Scan(context.TODO())
	if err != nil {
		m = nil
	}
	return
}
func (c *{{$.DaoName}}) GetByID(conn bun.IDB, id string) (m *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	m = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	m.ID = id
	err = conn.NewSelect().Model(m).WherePK().Scan(context.TODO())
	if err != nil {
		m = nil
	}
	return
}

//update before
//you must get OptimisticLockVersion value
func (c *{{$.DaoName}}) GetLockByID(conn bun.IDB, ids ...string) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	for i,_ := range ids {
		m := &{{$.ModelPkgName}}.{{$.TypeName}}{}
		m.ID = ids[i]
		ms = append(ms,m)
	}
	err = conn.NewSelect().Model(&ms).WherePK().Column({{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion,{{$.ModelPkgName}}.{{$.TypeName}}_Struct+"."+{{$.ModelPkgName}}.{{$.TypeName}}_ID).For("UPDATE").Scan(context.TODO())
	if err != nil {
		ms = nil
	}
	return
}
func (c *{{$.DaoName}}) GetLockByModelID(conn bun.IDB, ms ...*{{$.ModelPkgName}}.{{$.TypeName}}) error {
	return conn.NewSelect().Model(&ms).WherePK().Column({{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion,{{$.ModelPkgName}}.{{$.TypeName}}_Struct+"."+{{$.ModelPkgName}}.{{$.TypeName}}_ID).For("UPDATE").Scan(context.TODO())
}

func (c *{{$.DaoName}}) QueryWithLock(conn bun.IDB, condition string, params ...interface{}) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	if len(condition) < 1 {
		err = conn.NewSelect().Model(&ms).For("UPDATE").Scan(context.TODO())
	} else {
		err = conn.NewSelect().Model(&ms).Where(condition, params...).For("UPDATE").Scan(context.TODO())
	}
	if err != nil { //be sure
		ms = nil
	}
	return
}
func (c *{{$.DaoName}}) Query(conn bun.IDB, condition string, params ...interface{}) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	if len(condition) < 1 {
		err = conn.NewSelect().Model(&ms).Scan(context.TODO())
	} else {
		err = conn.NewSelect().Model(&ms).Where(condition, params...).Scan(context.TODO())
	}
	if err != nil { //be sure
		ms = nil
	}
	return
}

func (c *{{$.DaoName}}) ListWithLock(conn bun.IDB) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	err = conn.NewSelect().Model(&ms).For("UPDATE").Scan(context.TODO())
	if err != nil {//be sure
		ms = nil
	}
	return
}
func (c *{{$.DaoName}}) List(conn bun.IDB) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	err = conn.NewSelect().Model(&ms).Scan(context.TODO())
	if err != nil {//be sure
		ms = nil
	}
	return
}

func (c *{{$.DaoName}}) Count(conn bun.IDB, condition string, params ...interface{}) (count int, err error) {
	if len(condition) < 1 {
		count, err = conn.NewSelect().Model(&{{$.ModelPkgName}}.{{$.TypeName}}{}).Count(context.TODO())
	} else {
		count, err = conn.NewSelect().Model(&{{$.ModelPkgName}}.{{$.TypeName}}{}).Where(condition, params...).Count(context.TODO())
	}
	return
}

func (c *{{$.DaoName}}) QueryPageWithLock(conn bun.IDB, pageSize int, page int, condition string, params ...interface{}) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	if len(condition) < 1 {
		err = conn.NewSelect().Model(&ms).Limit(pageSize).Offset((page - 1) * pageSize).For("UPDATE").Scan(context.TODO())
	}else {
		err = conn.NewSelect().Model(&ms).Where(condition, params...).Limit(pageSize).Offset((page - 1) * pageSize).For("UPDATE").Scan(context.TODO())
	}
	if err != nil { //be sure
		ms = nil
	}
	return
}
func (c *{{$.DaoName}}) QueryPage(conn bun.IDB, pageSize int, page int, condition string, params ...interface{}) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	if len(condition) < 1 {
		err = conn.NewSelect().Model(&ms).Limit(pageSize).Offset((page - 1) * pageSize).Scan(context.TODO())
	}else {
		err = conn.NewSelect().Model(&ms).Where(condition, params...).Limit(pageSize).Offset((page - 1) * pageSize).Scan(context.TODO())
	}
	if err != nil { //be sure
		ms = nil
	}
	return
}

// count counts valid records which after conditions filter, rather than whole table's count
func (c *{{$.DaoName}}) QueryPageWithCount(
	conn bun.IDB,
	pageSize,
	pageNum int,
	condition string,
	params ...interface{},
) (ms []*{{$.ModelPkgName}}.{{$.TypeName}}, count int, err error) {
	if len(condition) < 1 {
		count, err = conn.NewSelect().Model(&ms).Limit(pageSize).Offset((pageNum - 1) * pageSize).ScanAndCount(context.TODO())
	} else {
		count, err = conn.NewSelect().Model(&ms).Where(condition, params...).Limit(pageSize).Offset((pageNum - 1) * pageSize).ScanAndCount(context.TODO())
	}

	if err != nil { //be sure
		ms = nil
	}
	return
}

func (c *{{$.DaoName}}) QueryOneWithLock(conn bun.IDB, condition string, params ...interface{}) (m *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	m = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	if len(condition) < 1 {
		err = conn.NewSelect().Model(m).For("UPDATE").Limit(1).Scan(context.TODO())
	} else {
		err = conn.NewSelect().Model(m).Where(condition, params...).For("UPDATE").Limit(1).Scan(context.TODO())
	}
	if err != nil {//be sure
		m = nil
	}
	return
}

func (c *{{$.DaoName}}) QueryOne(conn bun.IDB, condition string, params ...interface{}) (m *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	m = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	if len(condition) < 1 {
		err = conn.NewSelect().Model(m).Scan(context.TODO())
	} else {
		err = conn.NewSelect().Model(m).Where(condition, params...).Scan(context.TODO())
	}
	if err != nil {//be sure
		m = nil
	}
	return
}

func (c *{{$.DaoName}}) Insert(conn bun.IDB, m *{{$.ModelPkgName}}.{{$.TypeName}}) (err error) {
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuidForDB()
	}
	m.CreateTime = time.Now().Unix()
	m.UpdateTime = m.CreateTime
	_, err = conn.NewInsert().Model(m).Exec(context.TODO())
	return
}

func (c *{{$.DaoName}}) InsertReturn(conn bun.IDB, m *{{$.ModelPkgName}}.{{$.TypeName}}) ( mnew *{{$.ModelPkgName}}.{{$.TypeName}}, err error) {
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuidForDB()
	}
	m.CreateTime = time.Now().Unix()
	m.UpdateTime = m.CreateTime

	mnew = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	_, err = conn.NewInsert().Model(m).Returning("*").Exec(context.TODO(), mnew)
	if err != nil{
		mnew = nil
	}
	return
}

func (c *{{$.DaoName}}) Upsert(conn bun.IDB, m *{{$.ModelPkgName}}.{{$.TypeName}}) (err error) {
    m.UpdateTime = time.Now().Unix()
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuidForDB()
		m.CreateTime = m.UpdateTime
	} else if m.CreateTime == 0 {
		m.CreateTime = m.UpdateTime
	}
	m.OptimisticLockVersion++
	om := conn.NewInsert().Model(m).On("CONFLICT (id) DO UPDATE").Where({{$.ModelPkgName}}.{{$.TypeName}}_Struct+"."+{{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion+" = ?",m.OptimisticLockVersion-1)
	for _, it := range m.ToUpsertSet() {
		om.Set(it)
	}
	res, err := om.Exec(context.TODO())
	if n, _ := res.RowsAffected(); n == 0 {
		newm, err := c.GetLockByID(conn, m.ID)
		if err != nil {
			return err
		}
		m.OptimisticLockVersion = newm[0].OptimisticLockVersion
		err = c.Update(conn, m)
	}
	return err
}

func (c *{{$.DaoName}}) UpsertReturn(conn bun.IDB, m *{{$.ModelPkgName}}.{{$.TypeName}}) ( mnew *{{$.ModelPkgName}}.{{$.TypeName}},err error) {
	m.UpdateTime = time.Now().Unix()
	if len(m.ID) < 1 {
		m.ID = uuid.GetUuidForDB()
		m.CreateTime = m.UpdateTime
	} else if m.CreateTime == 0 {
		m.CreateTime = m.UpdateTime
	}
	m.OptimisticLockVersion++
	om := conn.NewInsert().Model(m).On("CONFLICT (id) DO UPDATE").Where({{$.ModelPkgName}}.{{$.TypeName}}_Struct+"."+{{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion+" = ?",m.OptimisticLockVersion-1)
	for _, it := range m.ToUpsertSet() {
		om.Set(it)
	}
	mnew = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	res, err := om.Returning("*").Exec(context.TODO(), mnew)
	if n, _ := res.RowsAffected(); n == 0 {
		ms, err := c.GetLockByID(conn, m.ID)
		if err != nil {
			return nil, err
		}
		m.OptimisticLockVersion = ms[0].OptimisticLockVersion
		mnew, err = c.UpdateReturn(conn, m)
	}
	return
}

func (c *{{$.DaoName}}) Update(conn bun.IDB, m *{{$.ModelPkgName}}.{{$.TypeName}}) (err error) {
	m.UpdateTime = time.Now().Unix()
	m.OptimisticLockVersion++
	res, err := conn.NewUpdate().Model(m).Where({{$.ModelPkgName}}.{{$.TypeName}}_ID+" = ? and "+{{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion+" = ?", m.ID, m.OptimisticLockVersion-1).Exec(context.TODO())
	if n, _ := res.RowsAffected(); n == 0 {
		err = sql.ErrNoRows
	}
	return
}

func (c *{{$.DaoName}}) UpdateReturn(conn bun.IDB, m *{{$.ModelPkgName}}.{{$.TypeName}}) (mnew *{{$.ModelPkgName}}.{{$.TypeName}},  err error) {
	m.UpdateTime = time.Now().Unix()
	m.OptimisticLockVersion++
	mnew = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	res, err := conn.NewUpdate().Model(m).Where({{$.ModelPkgName}}.{{$.TypeName}}_ID+" = ? and "+{{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion+" = ?", m.ID, m.OptimisticLockVersion-1).Returning("*").Exec(context.TODO(), mnew)
	if err != nil {
		mnew = nil
	}
	if n, _ := res.RowsAffected(); n == 0 {
		err = sql.ErrNoRows
	}
	return
}

func (c *{{$.DaoName}}) Delete(conn bun.IDB, m *{{$.ModelPkgName}}.{{$.TypeName}}) (err error) {
	return c.DeleteByID(conn, m.ID)
}

func (c *{{$.DaoName}}) DeleteByID(conn *bun.DB, id string) (err error) {
	_, err = conn.NewDelete().Model((*{{$.ModelPkgName}}.{{$.TypeName}})(nil)).Where({{$.ModelPkgName}}.{{$.TypeName}}_ID+" = ?", id).Exec(context.TODO())
	return
}

func (c *{{$.DaoName}}) DeleteByIDs(conn bun.IDB, ids []string, oneMax int) (err error) {
	m := (*{{$.ModelPkgName}}.{{$.TypeName}})(nil)
	max := oneMax
	times := len(ids)/max;
	for i := 1; i < times; i++ {
		oneIDs := ids[(i-1) * max:i * max -1]
		_, err = conn.NewDelete().Model(m).Where({{$.ModelPkgName}}.{{$.TypeName}}_ID+" in (?)", bun.In(oneIDs)).Exec(context.TODO())
		if err != nil {
			return
		}
	}

	if max * times < len(ids) {
		oneIDs := ids[max * times:]
		_, err = conn.NewDelete().Model(m).Where({{$.ModelPkgName}}.{{$.TypeName}}_ID+" in (?)", bun.In(oneIDs)).Exec(context.TODO())
	}
	return 
}

func (c *{{$.DaoName}}) DeleteReturn(conn bun.IDB, m *{{$.ModelPkgName}}.{{$.TypeName}}) (mnew *{{$.ModelPkgName}}.{{$.TypeName}},err error) {
	mnew = &{{$.ModelPkgName}}.{{$.TypeName}}{}
	_, err = conn.NewDelete().Model(m).WherePK().Returning("*").Exec(context.TODO(), mnew)
	if err != nil {
		mnew = nil
	}
	return
}

//example,please edit it
//update designated column with Optimistic Lock
func (c *{{$.DaoName}}) Update{{$.TypeName}}SomeColumn(conn bun.IDB, ids []string,/*todo: update parameters*/) (err error) {

	ms, err := c.GetLockByID(conn, ids...)
	if err != nil {
		return
	}
	ctx := context.TODO()
	condition := {{$.ModelPkgName}}.{{$.TypeName}}_ID+" = ? and "+{{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion+" = ?"
	for i, _ := range ms {
		ms[i].UpdateTime = time.Now().Unix()
		ms[i].OptimisticLockVersion++
		_, err = conn.NewUpdate().Model(ms[i]).Where(condition, ms[i].ID, ms[i].OptimisticLockVersion-1).Column(/*{{$.ModelPkgName}}.{{$.TypeName}}_xx,*/ {{$.ModelPkgName}}.{{$.TypeName}}_OptimisticLockVersion, {{$.ModelPkgName}}.{{$.TypeName}}_UpdateTime).Exec(ctx)
		if err != nil {
			dot.Logger().Debugln(err.Error())
			return
		}
	}
	return
}
`
