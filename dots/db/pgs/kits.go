package pgs

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

func CreateSchema(db *pg.DB, ms []interface{}) error {
	for _, model := range ms {
		err := db.CreateTable(model, &orm.CreateTableOptions{Temp: false, IfNotExists: true})
		if err != nil {
			return err
		}
	}
	return nil
}

//
func SqlLikeEscapePostgres(param string) string {
	return SqlLikeEscape(param, []rune{'%', '_'}, '\\')
}

func SqlLikeEscape(param string, special []rune, escape rune) string {
	cs := []rune(param)
	newcs := make([]rune, 0, len(cs)+10)
	es := escape
	for _, c := range cs {
		for _, s := range special {
			if c == s {
				newcs = append(newcs, es)
				break
			}
		}
		newcs = append(newcs, c)
	}

	res := string(newcs)
	return res
}
