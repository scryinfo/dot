package gorms

import (
	"github.com/scryinfo/dot/dot"
	"gorm.io/gorm"
)

const DaoBaseTypeID = "06ce4450-161e-44aa-8695-527a97bf57c5"

type DaoBase struct {
	Wrapper *Gorms `dot:""`
	//todo add
}

//WithTx with transaction, if return err != nil then rollback, or commit the transaction
func (c *DaoBase) WithTx(task func(tx *gorm.DB) error) error {
	var err error
	if task != nil {
		err = c.Wrapper.GetDb().Transaction(task)
	}
	return err
}

//DaoBaseTypeLives
func DaoBaseTypeLives() []*dot.TypeLives {
	tl := &dot.TypeLives{
		Meta: dot.Metadata{TypeID: DaoBaseTypeID, NewDoter: func(conf []byte) (dot.Dot, error) {
			return &DaoBase{}, nil
		}},
		Lives: []dot.Live{
			{
				LiveID:    DaoBaseTypeID,
				RelyLives: map[string]dot.LiveID{"Wrapper": GormTypeID},
			},
		},
	}

	lives := []*dot.TypeLives{tl}
	lives = append(lives, GormsTypeLives()...)
	return lives
}

//GenerateDaoBaseByDb this func is for test
func GenerateDaoBaseByDb(db *Gorms) *DaoBase {
	base := &DaoBase{db}
	return base
}
