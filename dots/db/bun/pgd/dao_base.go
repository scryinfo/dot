package pgd

import (
	"context"
	"database/sql"

	"github.com/scryinfo/dot/dot"
	"github.com/uptrace/bun"
)

const (
	//DaoBaseTypeID type id
	DaoBaseTypeID = "4d6adee7-7c10-4471-8885-a589688bac93"
)

// DaoBase doa base
type DaoBase struct {
	Wrapper *ConnWrapper `dot:""`
}

func (c *DaoBase) getConn() *bun.Conn {
	conn, _ := c.Wrapper.db.Conn(context.TODO())
	return &conn
}

// WithTx with transaction, if return err != nil then rollback, or commit the transaction
func (c *DaoBase) WithTx(task func(conn *bun.Tx) error) error {
	var err error
	if task != nil {
		conn := c.getConn()
		defer conn.Close()
		tx, err := conn.BeginTx(context.TODO(), &sql.TxOptions{})
		if err == nil {
			defer func() {
				if err == nil {
					err = tx.Commit()
				} else {
					err = tx.Rollback()
				}
			}()
			err = task(&tx)
		}
	}
	return err
}

// WithNoTx no transaction
func (c *DaoBase) WithNoTx(task func(conn *bun.DB) error) error {
	// todo

	return nil
}

// DaoBaseTypeLives make all type lives
func DaoBaseTypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeID: DaoBaseTypeID, NewDoter: func(conf []byte) (dot.Dot, error) {
			return &DaoBase{}, nil
		}},
		Lives: []dot.Live{
			{
				LiveID:    DaoBaseTypeID,
				RelyLives: map[string]dot.LiveID{"Wrapper": ConnWrapperTypeID},
			},
		},
	}
	lives := ConnWrapperTypeLives()
	lives = append(lives, tl)
	return lives
}

// GenerateDaoBase this func is for test
func GenerateDaoBase(conf string) *DaoBase {
	wrapper := GenerateConnWrapper(conf)
	base := &DaoBase{wrapper}
	return base
}

// GenerateDaoBaseByDb this func is for test
func GenerateDaoBaseByDb(db *bun.DB) *DaoBase {
	wrapper := GenerateConnWrapperByDb(db)
	base := &DaoBase{wrapper}
	return base
}
