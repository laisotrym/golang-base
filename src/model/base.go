package model

import (
    "database/sql"
)

type IBase interface {
    TableName() string
    GetColumns() []string
    InsertColumns() []string
}

type Base struct {
    CreatedBy sql.NullInt64 `db:"created_by" json:"created_by"`
    CreatedAt sql.NullTime  `db:"created_at" json:"created_at"`
    UpdatedBy sql.NullInt64 `db:"updated_by" json:"updated_by"`
    UpdatedAt sql.NullTime  `db:"updated_at" json:"updated_at"`
}

func (b Base) TableName() string {
    panic("implement me")
}

func (b Base) GetColumns() []string {
    panic("implement me")
}

func (b Base) InsertColumns() []string {
    panic("implement me")
}
