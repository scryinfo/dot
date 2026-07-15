package db

import "github.com/jackc/pgx/v5/pgtype"

func PgTextToString(v pgtype.Text) string {
	if !v.Valid {
		return ""
	}
	return v.String
}
func StringtToPgText(v string) pgtype.Text {
	return pgtype.Text{
		Valid:  true,
		String: v,
	}
}
