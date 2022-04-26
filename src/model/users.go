package model

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"

	"safeweb.app/rpc/safeweb_lib"
)

// User is generated type for table 'users'
type User struct {
	Base
	Id           int64                    `db:"id" json:"id"`
	UserName     string                   `db:"username" json:"username"`
	Email        string                   `db:"email" json:"email"`
	Phone        sql.NullString           `db:"phone" json:"phone"`
	FullName     sql.NullString           `db:"full_name" json:"full_name"`
	Country      sql.NullString           `db:"country" json:"country"`
	TimeZone     sql.NullString           `db:"time_zone" json:"time_zone"`
	PasswordHash string                   `db:"password_hash" json:"password_hash"`
	Status       safeweb_lib.ActiveStatus `db:"status" json:"status"`
	IsSuper      bool                     `db:"is_super" json:"is_super"`
	IsAdmin      bool                     `db:"is_admin" json:"is_admin"`
	IsMember     bool                     `db:"is_member" json:"is_member"`
	Notif        string                   `db:"notif" json:"notif"`
	ConfirmToken sql.NullString           `db:"confirm_token" json:"confirm_token"`
	ConfirmTime  sql.NullString           `db:"confirm_time" json:"confirm_time"`
}

// ComparePassword checks if the provided password is correct or not
func (u *User) ComparePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err
}

// Clone returns a clone of this user
func (u *User) Clone() *User {
	return &User{
		UserName:     u.UserName,
		Email:        u.Email,
		Phone:        u.Phone,
		FullName:     u.FullName,
		Country:      u.Country,
		TimeZone:     u.TimeZone,
		PasswordHash: u.PasswordHash,
		Status:       u.Status,
		IsSuper:      u.IsSuper,
		IsAdmin:      u.IsAdmin,
		IsMember:     u.IsMember,
		Notif:        u.Notif,
		ConfirmToken: u.ConfirmToken,
		ConfirmTime:  u.ConfirmTime,
	}
}

// table 'users' columns list struct
type tblUserColumns struct {
	Id           string
	Username     string
	Email        string
	Phone        string
	FullName     string
	Country      string
	TimeZone     string
	PasswordHash string
	Status       string
	IsSuper      string
	IsAdmin      string
	IsMember     string
	Notif        string
	ConfirmToken string
	ConfirmTime  string
	CreatedBy    string
	CreatedAt    string
	UpdatedBy    string
	UpdatedAt    string
}

// table 'users' metadata struct
type tblUser struct {
	Name    string
	Columns tblUserColumns
}

// table 'users' metadata info
var tblUserDefine = tblUser{
	Columns: tblUserColumns{
		Id:           "id",
		Username:     "username",
		Email:        "email",
		Phone:        "phone",
		FullName:     "full_name",
		Country:      "country",
		TimeZone:     "time_zone",
		Status:       "status",
		IsSuper:      "is_super",
		IsAdmin:      "is_admin",
		IsMember:     "is_member",
		PasswordHash: "password_hash",
		Notif:        "notif",
		ConfirmToken: "confirm_token",
		ConfirmTime:  "confirm_time",
		UpdatedAt:    "updated_at",
		UpdatedBy:    "updated_by",
		CreatedAt:    "created_at",
		CreatedBy:    "created_by",
	},
	Name: "core_user",
}

// InsertColumns return list columns name for table 'users'
func (*tblUser) GetColumns() []string {
	return []string{"id", "username", "email", "phone", "full_name", "country", "time_zone", "password_hash", "status", "is_super", "is_admin", "is_member", "notif", "confirm_token", "confirm_time", "created_by", "created_at", "updated_by", "updated_at"}
}

// InsertColumns return list columns name for table 'users'
func (*tblUser) InsertColumns() []string {
	return []string{"username", "email", "phone", "full_name", "country", "time_zone", "password_hash", "status", "is_super", "is_admin", "is_member", "notif", "confirm_token", "confirm_time", "created_by", "updated_by"}
}

// T return metadata info for table 'users'
func (u *User) T() *tblUser {
	return &tblUserDefine
}

// TableName return table name
func (u User) TableName() string {
	return "core_user"
}
