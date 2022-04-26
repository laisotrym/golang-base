package model

type Stock struct {
	Base
	Id         int64   `db:"id" json:"id"`
	MaCk       string  `db:"ma_ck" json:"ma_ck"`
	GiaMua     int64   `db:"gia_mua" json:"gia_mua"`
	KhoiLuong  int64   `db:"khoi_luong" json:"khoi_luong"`
	SoNgay     int64   `db:"so_ngay" json:"so_ngay"`
	TienLo     int64   `db:"tien_lo" json:"tien_lo"`
	TyLeLo     float32 `db:"ty_le_lo" json:"ty_le_lo"`
	GiaMax     int64   `db:"gia_ban_max" json:"gia_ban_max"`
	LaiMax     int64   `db:"lai_max" json:"lai_max"`
	TyLeMax    float32 `db:"ty_le_max" json:"ty_le_max"`
	GiaHomNay  int64   `db:"gia_ban_hom_nay" json:"gia_ban_hom_nay"`
	LaiHomNay  int64   `db:"lai_hom_nay" json:"lai_hom_nay"`
	TyLeHomNay float32 `db:"ty_le_hom_nay" json:"ty_le_hom_nay"`
	TrangThai  string  `db:"status" json:"trang_thai"`
}

// table 'stock' columns list struct
type tblStockColumns struct {
	Id         string
	MaCk       string
	GiaMua     string
	KhoiLuong  string
	SoNgay     string
	TienLo     string
	TyLeLo     string
	GiaMax     string
	LaiMax     string
	TyLeMax    string
	GiaHomNay  string
	LaiHomNay  string
	TyLeHomNay string
	TrangThai  string
	CreatedBy  string
	CreatedAt  string
	UpdatedBy  string
	UpdatedAt  string
}

// table 'stock' metadata struct
type tblStock struct {
	Name    string
	Columns tblStockColumns
}

// table 'users' metadata info
var tblStockDefine = tblStock{
	Columns: tblStockColumns{
		Id:         "id",
		MaCk:       "ma_ck",
		GiaMua:     "gia_mua",
		KhoiLuong:  "khoi_luong",
		SoNgay:     "so_ngay",
		TienLo:     "tien_lo",
		TyLeLo:     "ty_le_lo",
		GiaMax:     "gia_ban_max",
		LaiMax:     "lai_max",
		TyLeMax:    "ty_le_max",
		GiaHomNay:  "gia_ban_hom_nay",
		LaiHomNay:  "lai_hom_nay",
		TyLeHomNay: "ty_le_hom_nay",
		TrangThai:  "status",
		UpdatedAt:  "updated_at",
		UpdatedBy:  "updated_by",
		CreatedAt:  "created_at",
		CreatedBy:  "created_by",
	},
	Name: "laiit",
}

// InsertColumns return list columns name for table 'users'
func (*Stock) GetColumns() []string {
	return []string{"id", "ma_ck", "gia_mua", "khoi_luong",
		"so_ngay", "tien_lo", "ty_le_lo",
		"gia_ban_max", "lai_max", "ty_le_max",
		"gia_ban_hom_nay", "lai_hom_nay", "ty_le_hom_nay",
		"status",
		"created_by", "created_at", "updated_by", "updated_at"}
}

// InsertColumns return list columns name for table 'users'
func (*Stock) InsertColumns() []string {
	return []string{"ma_ck", "gia_mua", "khoi_luong",
		"so_ngay", "tien_lo", "ty_le_lo",
		"gia_ban_max", "lai_max", "ty_le_max",
		"gia_ban_hom_nay", "lai_hom_nay", "ty_le_hom_nay",
		"status",
		"created_by", "updated_by"}
}

// T return metadata info for table 'users'
func (u *Stock) T() *tblStock {
	return &tblStockDefine
}

// TableName return table name
func (u Stock) TableName() string {
	return "laiit"
}
