package pgs

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"reflect"
	"strconv"
	"strings"
)

//CreateSchema create tables
func CreateSchema(db *pg.DB, ms []interface{}) error {
	for _, model := range ms {
		err := db.CreateTable(model, &orm.CreateTableOptions{Temp: false, IfNotExists: true})
		if err != nil {
			return err
		}
	}
	return nil
}

//SQLLikeEscapePostgres like escape
func SQLLikeEscapePostgres(param string) string {
	return SQLLikeEscape(param, []rune{'%', '_'}, '\\')
}

//SQLLikeEscape  like escape
func SQLLikeEscape(param string, special []rune, escape rune) string {
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

//ToMap convert fields of struct to map,  key is Lowercase for the first letter of the field name
func ToMap(d interface{}, ex map[string]bool) map[string]string {
	if ex == nil {
		ex = map[string]bool{}
	}
	res := make(map[string]string)
	mv := reflect.ValueOf(d).Elem()
	mt := mv.Type()
	for i := 0; i < mv.NumField(); i++ {
		v := mv.Field(i)
		fname := mt.Field(i).Name
		lfname := ToFirstLower(fname)
		if _, ok := ex[fname]; ok {
			continue
		}

		if _, ok := ex[lfname]; ok {
			continue
		}
		switch v.Kind() {
		case reflect.Int64, reflect.Int32, reflect.Int, reflect.Int16:
			res[lfname] = strconv.FormatInt(v.Int(), 10)
		default:
			res[lfname] = v.String()
		}
	}

	return res
}

//ToFirstLower Lowercase for the first letter
func ToFirstLower(s string) string {
	cs := []rune(s)
	return strings.ToLower(string(cs[0])) + string(cs[1:])
}

/////////////////////////////////////////////see github.com/go-pg/pg/internal
func isUpper(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func isLower(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func toUpper(c byte) byte {
	return c - 32
}

func toLower(c byte) byte {
	return c + 32
}

// Underscore converts "CamelCasedString" to "camel_cased_string".
// Just work for ASCⅡ code
func Underscore(s string) string {
	r := make([]byte, 0, len(s)+5)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if isUpper(c) {
			if i > 0 && i+1 < len(s) && (isLower(s[i-1]) || isLower(s[i+1])) {
				r = append(r, '_', toLower(c))
			} else {
				r = append(r, toLower(c))
			}
		} else {
			r = append(r, c)
		}
	}
	return string(r)
}

//CamelCased Just work for ASCⅡ code
func CamelCased(s string) string {
	r := make([]byte, 0, len(s))
	upperNext := true
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' {
			upperNext = true
			continue
		}
		if upperNext {
			if isLower(c) {
				c = toUpper(c)
			}
			upperNext = false
		}
		r = append(r, c)
	}
	return string(r)
}

//ToExported Just work for ASCⅡ code
func ToExported(s string) string {
	if len(s) == 0 {
		return s
	}
	if c := s[0]; isLower(c) {
		b := []byte(s)
		b[0] = toUpper(c)
		return string(b)
	}
	return s
}

//UpperString Just work for ASCⅡ code
func UpperString(s string) string {
	if isUpperString(s) {
		return s
	}

	b := make([]byte, len(s))
	for i := range b {
		c := s[i]
		if isLower(c) {
			c = toUpper(c)
		}
		b[i] = c
	}
	return string(b)
}

func isUpperString(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if isLower(c) {
			return false
		}
	}
	return true
}

/////////////////////////////////////////////see github.com/go-pg/pg/internal
