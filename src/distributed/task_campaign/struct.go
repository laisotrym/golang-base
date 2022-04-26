package task_campaign

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"safeweb.app/model/enum"
	safeweb_lib_num "safeweb.app/service/safeweb_lib/enum"
)safeweb.app

type ZData struct {
	Score  int64
	Member interface{}
}

// =========================== Define BaseInfo ===========================//
type BaseInfo struct {
	CampaignId     int64                        `json:"campaign_id" yaml:"campaign_id"`
	CampaignName   string                       `json:"campaign_name" yaml:"campaign_name"`
	VoucherCode    string                       `json:"voucher_code" yaml:"voucher_code"`
	FromDate       *time.Time                   `json:"from_date" yaml:"from_date"`
	ToDate         *sql.NullTime                `json:"to_date" yaml:"to_date"`
	QrTypeList     []safeweb_lib_num.QrType     `json:"qr_type_list" yaml:"qr_type_list"`
	TypeSourceList []safeweb_lib_num.TypeSource `json:"type_source_list" yaml:"type_source_list"`
	PayMethodList  []safeweb_lib_num.PayMethod  `json:"pay_method_list" yaml:"pay_method_list"`
}

// ========================== End For BaseInfo ========================== //
// =========================== Define ITaskCampaign ===========================//
type ITaskCampaign interface {
	GetCampaignId() int64
	VoucherCodeColumns(ITaskCampaign) []string
	Columns(ITaskCampaign) []string
	GroupColumns(ITaskCampaign) []string
	WhereConditions(ITaskCampaign) string
	WhereArgs() []interface{}
	TableName() string
	KeyRedis() string
	FieldRedis() *string
	DataRedis() interface{}
	DataRedisSortedSetsScore(int64) int64
	DataRedisSortedSets() (interface{}, error)
	GetFromDate() *time.Time
	GetToDate() *sql.NullTime
	GetCampFromDate() *time.Time
	GetCampToDate() *time.Time
}

// ========================== End For ITaskCampaign ========================== //
type CampaignRange struct {
	CampFromDate *time.Time `db:"camp_from_date" json:"camp_from_date" yaml:"camp_from_date"`
	CampToDate   *time.Time `db:"camp_to_date" json:"camp_to_date" yaml:"camp_to_date"`
}

func (t CampaignRange) GetCampFromDate() *time.Time {
	return t.CampFromDate
}

func (t CampaignRange) GetCampToDate() *time.Time {
	return t.CampToDate
}

// =========================== Define DateRange ===========================//
type DateRange struct {
	FromDate *time.Time    `db:"from_date" json:"from_date" yaml:"from_date"`
	ToDate   *sql.NullTime `db:"to_date" json:"to_date" yaml:"to_date"`
	CampaignRange
}

func (t DateRange) GetCampaignId() int64 {
	panic("implement me")
}

func (t DateRange) VoucherCodeColumns(obj ITaskCampaign) []string {
	panic("implement me")
}

func (t DateRange) Columns(obj ITaskCampaign) []string {
	panic("implement me")
}

func (t DateRange) WhereConditions(obj ITaskCampaign) string {
	panic("implement me")
}

func (t DateRange) WhereArgs() []interface{} {
	panic("implement me")
}

func (t DateRange) GroupColumns(obj ITaskCampaign) []string {
	panic("implement me")
}

func (t DateRange) TableName() string {
	panic("implement me")
}

func (t DateRange) KeyRedis() string {
	panic("implement me")
}

func (t DateRange) FieldRedis() *string {
	panic("implement me")
}

func (t DateRange) DataRedis() interface{} {
	panic("implement me")
}

func (t DateRange) DataRedisSortedSetsScore(in int64) int64 {
	return in
}

func (t DateRange) DataRedisSortedSets() (interface{}, error) {
	panic("implement me")
}

func (t DateRange) GetFromDate() *time.Time {
	return t.FromDate
}

func (t DateRange) GetToDate() *sql.NullTime {
	return t.ToDate
}

// ========================== End Fot DateRange ========================== //
// =========================== Define TaskCampaign ===========================//
type TaskCampaign struct {
	CampaignId int64 `db:"campaign_id" json:"campaign_id" yaml:"campaign_id"`
	DateRange
}

func (t TaskCampaign) VoucherCodeColumns(obj ITaskCampaign) []string {
	return append(obj.Columns(obj), []string{
		"min(campaign_vouchers.from_date) as camp_from_date",
		"max(campaign_vouchers.to_date) as camp_to_date",
	}...)
}

func (t TaskCampaign) Columns(obj ITaskCampaign) []string {
	tblName := obj.TableName()
	return append(obj.GroupColumns(obj), []string{
		fmt.Sprintf("min(%s.from_date) as from_date", tblName),
		fmt.Sprintf("max(%s.to_date) as to_date", tblName),
	}...)
}

func (t TaskCampaign) GetCampaignId() int64 {
	return t.CampaignId
}

func (t TaskCampaign) WhereConditions(obj ITaskCampaign) string {
	tblName := obj.TableName()
	return fmt.Sprintf(
		"%s.campaign_id = ? AND %s.deleted_at IS NULL AND (%s.to_date IS NULL OR %s.to_date >= now())",
		tblName, tblName, tblName, tblName,
	)
}

func (t TaskCampaign) WhereArgs() []interface{} {
	return []interface{}{
		t.CampaignId,
	}
}

func (t TaskCampaign) DataRedis() interface{} {
	return t.DateRange
}

// ========================== End For TaskCampaign ========================== //
// =========================== Define CampaignVouchers ===========================//
type CampaignVouchers struct {
	TaskCampaign
	VoucherCode string `db:"voucher_code" json:"voucher_code" yaml:"voucher_code"`
}

func (t CampaignVouchers) GroupColumns(obj ITaskCampaign) []string {
	tblName := obj.TableName()
	return []string{
		fmt.Sprintf("%s.campaign_id", tblName),
		fmt.Sprintf("%s.voucher_code", tblName),
	}
}

func (t CampaignVouchers) TableName() string {
	return "campaign_vouchers"
}

func (t CampaignVouchers) KeyRedis() string {
	return fmt.Sprintf("%s", t.VoucherCode)
}

func (t CampaignVouchers) FieldRedis() *string {
	return &t.VoucherCode
}

// ========================== End For CampaignVouchers========================== //
// =========================== Define CampaignMerchant ===========================//
type CampaignMerchant struct {
	TaskCampaign
	MerchantCode string `db:"merchant_code" json:"merchant_code" yaml:"merchant_code"`
}

func (t CampaignMerchant) GroupColumns(obj ITaskCampaign) []string {
	tblName := obj.TableName()
	return []string{
		fmt.Sprintf("%s.campaign_id", tblName),
		fmt.Sprintf("%s.merchant_code", tblName),
	}
}

// Return table name for table 'campaign_terminals' using for task_campaign.GetMerchant
func (t CampaignMerchant) TableName() string {
	return "campaign_terminals"
}

// Return redis key store using for task_campaign.GetMerchant
func (t CampaignMerchant) KeyRedis() string {
	return fmt.Sprintf(safeweb_lib_constant.FormatRedisMerchantKey, t.CampaignId)
}

func (t CampaignMerchant) FieldRedis() *string {
	return &t.MerchantCode
}

func (t CampaignMerchant) DataRedisSortedSets() (interface{}, error) {
	if obj := t.FieldRedis(); obj != nil {
		return *t.FieldRedis(), nil
	} else {
		return nil, nil
	}
}

// ========================== End For CampaignMerchant ========================== //
// =========================== Define CampaignTerminal ===========================//
type CampaignTerminal struct {
	CampaignMerchant
	TerminalCode string `db:"terminal_code" json:"terminal_code" yaml:"terminal_code"`
	TerminalId   string `db:"terminal_id" json:"terminal_id" json:"terminal_id" yaml:"terminal_id"`
}

func (t CampaignTerminal) GroupColumns(obj ITaskCampaign) []string {
	tblName := obj.TableName()
	return []string{
		fmt.Sprintf("%s.campaign_id", tblName),
		fmt.Sprintf("%s.merchant_code", tblName),
		fmt.Sprintf("%s.terminal_id", tblName),
	}
}

// Return table name for table 'campaign_terminals' using for task_campaign.GetTerminal
func (t CampaignTerminal) TableName() string {
	return "campaign_terminals"
}

// Return redis key store using for task_campaign.GetTerminal
func (t CampaignTerminal) KeyRedis() string {
	return fmt.Sprintf(safeweb_lib_constant.FormatRedisTerminalKey, t.CampaignId)
}

func (t CampaignTerminal) FieldRedis() *string {
	if len(t.TerminalCode) > 0 && t.TerminalCode != safeweb_lib_constant.StarStr {
		data := fmt.Sprintf("%s.%s", t.MerchantCode, t.TerminalCode)
		return &data
	} else if len(t.TerminalId) > 0 && t.TerminalId != safeweb_lib_constant.StarStr {
		data := fmt.Sprintf("%s.%s", t.MerchantCode, t.TerminalId)
		return &data
	} else {
		return nil
	}
}

func (t CampaignTerminal) DataRedisSortedSets() (interface{}, error) {
	if obj := t.FieldRedis(); obj != nil {
		return *t.FieldRedis(), nil
	} else {
		return nil, nil
	}
}

// ========================== End For CampaignTerminal========================== //
// =========================== Define CampaignBenefit ===========================//
type CampaignBenefit struct {
	TaskCampaign
	Type              enum.BenefitType `db:"type" json:"type" yaml:"type"`
	TnxMinAmount      float64          `db:"tnx_min_amount" json:"tnx_min_amount" yaml:"tnx_min_amount"`
	DiscountPercent   sql.NullFloat64  `db:"discount_percent" json:"discount_percent" yaml:"discount_percent"`
	MaxDiscountAmount float64          `db:"max_discount_amount" json:"max_discount_amount" yaml:"max_discount_amount"`
}

func (t CampaignBenefit) GroupColumns(obj ITaskCampaign) []string {
	tblName := obj.TableName()
	return []string{
		fmt.Sprintf("%s.campaign_id", tblName),
		fmt.Sprintf("%s.type", tblName),
		fmt.Sprintf("%s.tnx_min_amount", tblName),
		fmt.Sprintf("%s.discount_percent", tblName),
		fmt.Sprintf("%s.max_discount_amount", tblName),
	}
}

// Return table name for table 'campaign_benefits' using for task_campaign.GetBenefit
func (t CampaignBenefit) TableName() string {
	return "campaign_benefits"
}

// Return redis key store using for task_campaign.GetBenefit
func (t CampaignBenefit) KeyRedis() string {
	return fmt.Sprintf(safeweb_lib_constant.FormatRedisBenefitKey, t.CampaignId)
}

func (t CampaignBenefit) FieldRedis() *string {
	return nil
}

func (t CampaignBenefit) DataRedis() interface{} {
	return t.DateRange
}

func (t CampaignBenefit) DataRedisSortedSets() (interface{}, error) {
	discountPercent, _ := t.DiscountPercent.Value()
	return json.Marshal(map[string]interface{}{
		"tnx_min_amount":      t.TnxMinAmount,
		"discount_percent":    discountPercent,
		"max_discount_amount": t.MaxDiscountAmount,
	})
}

func (t CampaignBenefit) DataRedisSortedSetsScore(in int64) int64 {
	return in + int64(t.Type)
}

// ========================== End For CampaignBenefit ========================== //
// ========================== Define For CampaignBank========================== //
type CampaignBank struct {
	TaskCampaign
	BankCode string `db:"bank_code" json:"bank_code" yaml:"bank_code"`
}

func (t CampaignBank) GroupColumns(obj ITaskCampaign) []string {
	tblName := obj.TableName()
	return []string{
		fmt.Sprintf("%s.campaign_id", tblName),
		fmt.Sprintf("%s.bank_code", tblName),
	}
}

// Return table name for table 'campaign_banks' using for task_campaign.GetBanks
func (t CampaignBank) TableName() string {
	return "campaign_banks"
}

// Return redis key store using for task_campaign.GetBanks
func (t CampaignBank) KeyRedis() string {
	return fmt.Sprintf(safeweb_lib_constant.FormatRedisBankKey, t.CampaignId)
}

func (t CampaignBank) FieldRedis() *string {
	return &t.BankCode
}

func (t CampaignBank) DataRedisSortedSets() (interface{}, error) {
	if obj := t.FieldRedis(); obj != nil {
		return *t.FieldRedis(), nil
	} else {
		return nil, nil
	}
}

// ========================== End For CampaignBank ========================== //
// ========================== Start For CampaignDetail ========================== //
type CampaignDetail struct {
	TaskCampaign
}

func (t CampaignDetail) TableName() string {
	return "campaign_detail"
}

func (t CampaignDetail) Columns(obj ITaskCampaign) []string {
	return []string{safeweb_lib_constant.StarStr}
}

// ========================== End For CampaignDetail ========================== //
// ========================== Start For Campaign ========================== //
type Campaign struct {
	TaskCampaign
}

func (t Campaign) TableName() string {
	return "campaigns"
}

func (t Campaign) Columns(obj ITaskCampaign) []string {
	return []string{safeweb_lib_constant.StarStr}
}

// ========================== End For CampaignDetail ========================== //
// ========================== Start For BlackListCustomer ========================== //
type BlackListCustomer struct {
	DateRange
	Mobile string `db:"mobile" json:"phone" yaml:"mobile"`
}

func (t BlackListCustomer) Columns(obj ITaskCampaign) []string {
	tblName := obj.TableName()
	return append(obj.GroupColumns(obj), []string{
		fmt.Sprintf("min(%s.from_date) as from_date", tblName),
		fmt.Sprintf("max(%s.to_date) as to_date", tblName),
	}...)
}

func (t BlackListCustomer) GroupColumns(obj ITaskCampaign) []string {
	tblName := obj.TableName()
	return []string{
		fmt.Sprintf("%s.mobile", tblName),
	}
}

func (t BlackListCustomer) TableName() string {
	return "blacklist_customers"
}

func (t BlackListCustomer) KeyRedis() string {
	return fmt.Sprintf("%s", t.Mobile)
}

func (t BlackListCustomer) FieldRedis() *string {
	return &t.Mobile
}

func (t BlackListCustomer) DataRedis() interface{} {
	return t.DateRange
}

func (t BlackListCustomer) DataRedisSortedSets() (interface{}, error) {
	if obj := t.FieldRedis(); obj != nil {
		return *t.FieldRedis(), nil
	} else {
		return nil, nil
	}
}

func (t BlackListCustomer) WhereConditions(obj ITaskCampaign) string {
	tblName := obj.TableName()
	return fmt.Sprintf(
		"%s.mobile = ? AND %s.deleted_at IS NULL AND (%s.to_date IS NULL OR %s.to_date >= now())",
		tblName, tblName, tblName, tblName,
	)
}

func (t BlackListCustomer) WhereArgs() []interface{} {
	return []interface{}{
		t.Mobile,
	}
}

// ========================== End For BlackListCustomer ========================== //

// ========================== Start For BlackListTerminal ========================== //
type BlackListTerminal struct {
	DateRange
	MerchantCode string `db:"merchant_code" json:"merchant_code" yaml:"merchant_code"`
	TerminalCode string `db:"terminal_id" json:"terminal_id" yaml:"terminal_id"`
}

func (t BlackListTerminal) Columns(obj ITaskCampaign) []string {
	tblName := obj.TableName()
	return append(obj.GroupColumns(obj), []string{
		fmt.Sprintf("min(%s.from_date) as from_date", tblName),
		fmt.Sprintf("max(%s.to_date) as to_date", tblName),
	}...)
}

func (t BlackListTerminal) GroupColumns(obj ITaskCampaign) []string {
	tblName := obj.TableName()
	return []string{
		fmt.Sprintf("%s.merchant_code", tblName),
		fmt.Sprintf("%s.terminal_id", tblName),
	}
}

func (t BlackListTerminal) TableName() string {
	return "blacklist_terminals"
}

func (t BlackListTerminal) KeyRedis() string {
	return fmt.Sprintf(safeweb_lib_constant.FormatRedisBlTerminalKey, t.MerchantCode)
}

func (t BlackListTerminal) FieldRedis() *string {
	return &t.TerminalCode
}

func (t BlackListTerminal) DataRedis() interface{} {
	return t.DateRange
}

func (t BlackListTerminal) DataRedisSortedSets() (interface{}, error) {
	if obj := t.FieldRedis(); obj != nil {
		return *t.FieldRedis(), nil
	}
	return nil, nil
}

func (t BlackListTerminal) WhereConditions(obj ITaskCampaign) string {
	tblName := obj.TableName()
	return fmt.Sprintf(
		"%s.merchant_code = ? AND %s.deleted_at IS NULL AND (%s.to_date IS NULL OR %s.to_date >= now())",
		tblName, tblName, tblName, tblName,
	)
}

func (t BlackListTerminal) WhereArgs() []interface{} {
	return []interface{}{
		t.MerchantCode,
	}
}

// ========================== End For BlackListCustomer ========================== //
