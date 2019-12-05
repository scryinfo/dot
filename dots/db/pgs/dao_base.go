package pgs

import (
	"github.com/scryinfo/dot/dot"
)

const (
	DaoBaseTypeId = "4d6adee7-7c10-4471-8885-a589688bac93"
)

type DaoBase struct {
	Wrapper *ConnWrapper `dot:""`
}

func (c *DaoBase) getConn() *pg.Conn {
	conn := c.Wrapper.db.Conn()
	return conn
}

func (c *DaoBase) WithTx(tast func(conn *pg.Conn) error) error {
	var err error
	if tast != nil {
		conn := c.getConn()
		defer conn.Close()
		var tx *pg.Tx
		tx, err = conn.Begin()
		if err == nil {
			err = tast(conn)
			if err == nil {
				err = tx.Commit()
			} else {
				err = tx.Rollback()
			}
		}
	}
	return err
}

func (c *DaoBase) WithNoTx(tast func(conn *pg.Conn) error) error {
	var err error
	if tast != nil {
		conn := c.getConn()
		defer conn.Close()
		err = tast(conn)
	}
	return err
}

//DaoBaseTypeLives make all type lives
func DaoBaseTypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeId: DaoBaseTypeId, NewDoter: func(conf []byte) (dot.Dot, error) {
			return &DaoBase{}, nil
		}},
		Lives: []dot.Live{
			dot.Live{
				LiveId:    DaoBaseTypeId,
				RelyLives: map[string]dot.LiveId{"Wrapper": ConnWrapperTypeId},
			},
		},
	}

	lives := ConnWrapperTypeLives()
	lives = append(lives, tl)

	return lives
}

//GenerateDaoBase this func is for test
func GenerateDaoBase(conf string) *DaoBase {
	wrapper := GenarateConnWrapper(conf)
	base := &DaoBase{wrapper}
	return base
}
