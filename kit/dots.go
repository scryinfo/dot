package kit

import uuid "github.com/satori/go.uuid"

//GenerateUuid 生成uuid
func GenerateUuid() string {
	u1 := uuid.Must(uuid.NewV4())
	return u1.String()
}
