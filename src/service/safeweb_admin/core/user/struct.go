package user

import (
	"database/sql"
	"time"

	"safeweb.app/rpc/safeweb_lib"
)

type (
	ItemUser struct {
		Id           int64
		UserName     string
		Email        string
		Phone        sql.NullString
		FullName     sql.NullString
		Country      sql.NullString
		TimeZone     sql.NullString
		Password     string
		Status       safeweb_lib.ActiveStatus
		IsSuper      bool
		IsAdmin      bool
		IsMember     bool
		Notif        sql.NullString
		ConfirmToken string
		ConfirmTime  time.Time
		CreatedBy    int64
		CreatedAt    int64
		UpdatedBy    int64
		UpdatedAt    int64
	}
	CreateInput struct {
		User ItemUser
	}
	CreateOut struct {
		User ItemUser
	}
	ReadInput struct {
		Id int64
	}
	ReadOut struct {
		User ItemUser
	}
	UpdateInput struct {
		User ItemUser
	}
	UpdateOut struct {
		User ItemUser
	}
	DeleteInput struct {
		Id int64
	}
	DeleteOut struct {
		Success bool
	}
	SearchInput struct {
		Id       int64
		UserName string
		Email    string
		Phone    string
		FullName string
		Status   safeweb_lib.ActiveStatus
		IsSuper  bool
		IsAdmin  bool
		Page     int64 `validate:"gte=0"`
		PageSize int64 `validate:"gt=0"`
	}
	ListOutput struct {
		Users []*ItemUser
		Total int64
	}
)
