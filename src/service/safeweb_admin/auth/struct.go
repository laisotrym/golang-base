package auth

import (
	"database/sql"
)

type (
	ItemUser struct {
		Id       int64
		UserName string
		FullName sql.NullString
		Email    string
		Phone    sql.NullString
		Country  sql.NullString
		TimeZone sql.NullString
	}
)
