package buns

import (
	"context"

	"github.com/uptrace/bun"
)

// DaoBase doa base
type DaoBase struct {
	ConnWrapper *ConnWrapper
}

func (c *DaoBase) getConn() *bun.Conn {
	conn, _ := c.ConnWrapper.db.Conn(context.TODO())
	return &conn
}

// WithTx with transaction, if return err != nil then rollback, or commit the transaction
func (c *DaoBase) WithTx(task func(conn bun.IDB) error) error {
	var err error
	if task != nil {
		err = c.ConnWrapper.RunInTx(task)
	}
	return err
}

// WithNoTx no transaction
func (c *DaoBase) WithNoTx(task func(conn bun.IDB) error) error {
	var err error
	if task != nil {
		err = c.ConnWrapper.RunInNoTx(task)
	}

	return err
}

func NewDaoBase(connWrapper *ConnWrapper) *DaoBase {
	base := &DaoBase{connWrapper}
	return base
}
