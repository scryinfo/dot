package tools

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
