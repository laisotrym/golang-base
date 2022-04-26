package task_campaign_test

/*import (
    "context"
    "database/sql"
    "fmt"
    "regexp"
    "testing"
    "time"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"

    "safeweb.app/config"
    "safeweb.app/distributed/task_campaign"
    "safeweb.app/model"
    "safeweb.app/model/enum"
)

func TestRepository_StoreRedis(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    conn, _, red, err := config.InitMockConn(true)
    if err != nil {
        t.Fatal(err)
    }
    repo := task_campaign.NewRepository(conn, red)

    t.Run("ClientRuleMerchant", func(t *testing.T) {
        err = repo.StoreRedisHSet(ctx,
            task_campaign.CampaignMerchant{TaskCampaign: task_campaign.TaskCampaign{CampaignId: 1}},
            map[string]interface{}{"test": "1"},
            time.Now().Add(time.Duration(60*time.Second)).Unix(),
        )
        assert.NoError(t, err)
    })

    t.Run("Terminal", func(t *testing.T) {
        err = repo.StoreRedisHSet(ctx,
            task_campaign.CampaignTerminal{CampaignMerchant: task_campaign.CampaignMerchant{TaskCampaign: task_campaign.TaskCampaign{CampaignId: 1}}},
            map[string]interface{}{"test": "1"},
            time.Now().Add(time.Duration(60*time.Second)).Unix(),
        )
        assert.NoError(t, err)
    })

    t.Run("ClientRuleBank", func(t *testing.T) {
        err = repo.StoreRedisHSet(ctx,
            task_campaign.CampaignBank{TaskCampaign: task_campaign.TaskCampaign{CampaignId: 1}},
            map[string]interface{}{"VCB": "1"},
            time.Now().Add(time.Duration(60*time.Second)).Unix(),
        )
        assert.NoError(t, err)
    })

    t.Run("BaseInfo", func(t *testing.T) {
        var baseInfo = &task_campaign.BaseInfo{
            CampaignId:   1,
            CampaignName: "name",
        }
        err = repo.StoreRedisSet(ctx,
            task_campaign.CampaignVouchers{TaskCampaign: task_campaign.TaskCampaign{CampaignId: 1}, VoucherCode: "VOUCHER_VCB"},
            baseInfo, 20,
        )
        assert.NoError(t, err)
    })
}

func TestRepository_FindCampaign(t *testing.T) {
    conn, mock, red, err := config.InitMockConn(true)
    if err != nil {
        t.Fatal(err)
    }
    repo := task_campaign.NewRepository(conn, red)
    tests := []config.TestStruct{
        {
            Obj: model.Campaign{
                ID:     1,
                Name:   "Campaign VNPAY",
                Status: enum.Approved,
            },
            Query:       "SELECT * FROM campaigns WHERE (id = 1 and deleted_at IS NULL and status = 4)",
            EmptyError:  true,
            EmptyOutput: false,
        },
    }
    for i, test := range tests {
        t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
            if err != nil {
                t.Fatalf("Error  %v not expect when mock db", err)
            }

            item := test.Obj.(model.Campaign)
            var rows = sqlmock.NewRows([]string{"id", "name", "status"}).
                AddRow(item.ID, item.Name, item.Status)

            mock.ExpectQuery(regexp.QuoteMeta(test.Query)).WillReturnRows(rows)
            campaignObj, err := repo.GetCampaign(task_campaign.Campaign{TaskCampaign: task_campaign.TaskCampaign{CampaignId: item.ID}})

            assert.Equal(t, &item, campaignObj)
            assert.NoError(t, err)
            assert.NoError(t, mock.ExpectationsWereMet())
        })
    }
}

func TestRepository_FindCampaignDetail(t *testing.T) {
    conn, mock, red, err := config.InitMockConn(true)
    if err != nil {
        t.Fatal(err)
    }
    repo := task_campaign.NewRepository(conn, red)
    tests := []config.TestStruct{
        {
            Obj: model.CampaignDetail{
                CampaignID:        1,
                MaxBudgetAmount:   sql.NullFloat64{Float64: 250, Valid: true},
                MaxDiscountAmount: sql.NullFloat64{Float64: 10, Valid: true},
                UsageMaxCust:      10,
                UsagePeriodCust:   enum.Daily,
                PayMethodList:     "01, 02",
            },
            Query:       "SELECT * FROM campaign_detail WHERE (campaign_id = 1 and deleted_at IS NULL)",
            EmptyError:  true,
            EmptyOutput: false,
        },
    }
    for i, test := range tests {
        t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
            if err != nil {
                t.Fatalf("Error  %v not expect when mock db", err)
            }

            item := test.Obj.(model.CampaignDetail)
            var rows = sqlmock.NewRows([]string{"campaign_id", "max_discount_amount", "max_budget_amount", "usage_max_cust", "usage_period_cust", "pay_method_list",
                "created_by", "created_at", "updated_at", "updated_by", "deleted_by", "deleted_at"}).
                AddRow(item.CampaignID, item.MaxDiscountAmount, item.MaxBudgetAmount, item.UsageMaxCust, item.UsagePeriodCust, item.PayMethodList, item.CreatedBy,
                    item.CreatedAt, item.UpdatedAt, item.UpdatedBy, item.DeletedBy, item.DeletedAt)

            mock.ExpectQuery(regexp.QuoteMeta(test.Query)).WillReturnRows(rows)
            objResult, err := repo.GetCampaignDetail(task_campaign.CampaignDetail{TaskCampaign: task_campaign.TaskCampaign{CampaignId: item.CampaignID}})

            assert.Equal(t, &item, objResult)
            assert.NoError(t, err)
            assert.NoError(t, mock.ExpectationsWereMet())
        })
    }
}

func TestRepository_FindCampaignVouchers(t *testing.T) {
    conn, mock, red, err := config.InitMockConn(true)
    if err != nil {
        t.Fatal(err)
    }
    now := time.Now()
    repo := task_campaign.NewRepository(conn, red)
    campId := 1
    query := "SELECT campaign_vouchers.campaign_id, campaign_vouchers.voucher_code, min(campaign_vouchers.from_date) as from_date, max(campaign_vouchers.to_date) as to_date FROM campaign_vouchers WHERE (campaign_vouchers.campaign_id = (1) AND campaign_vouchers.deleted_at IS NULL AND (campaign_vouchers.to_date IS NULL OR campaign_vouchers.to_date >= now())) GROUP BY campaign_vouchers.campaign_id, campaign_vouchers.voucher_code"
    tests := []config.TestStruct{
        {
            Obj: task_campaign.CampaignVouchers{
                TaskCampaign: task_campaign.TaskCampaign{CampaignId: int64(campId),
                    DateRange: task_campaign.DateRange{
                        FromDate: &now,
                        ToDate:   nil,
                    },
                },
                VoucherCode: "VOUCHER_CODE_VCB",
            },
            Query:       query,
            EmptyError:  true,
            EmptyOutput: false,
        },
        {
            Obj: task_campaign.CampaignVouchers{
                TaskCampaign: task_campaign.TaskCampaign{CampaignId: int64(campId),
                    DateRange: task_campaign.DateRange{
                        FromDate: &now,
                        ToDate:   nil,
                    },
                },
                VoucherCode: "VOUCHER_CODE_VCB2",
            },
            Query:       query,
            EmptyError:  true,
            EmptyOutput: false,
        },
    }

    t.Run(fmt.Sprintf("CampVouchers"), func(t *testing.T) {
        if err != nil {
            t.Fatalf("Error  %v not expect when mock db", err)
        }

        item1 := tests[0].Obj.(task_campaign.CampaignVouchers)
        item2 := tests[1].Obj.(task_campaign.CampaignVouchers)
        var rows = sqlmock.NewRows([]string{"campaign_id", "voucher_code", "from_date", "to_date"}).
            AddRow(item1.TaskCampaign.CampaignId, item1.VoucherCode, item1.FromDate, item1.ToDate).
            AddRow(item2.TaskCampaign.CampaignId, item2.VoucherCode, item2.FromDate, item2.ToDate)

        mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
        count, campResults, err := repo.GetTaskCampaigns(task_campaign.CampaignVouchers{TaskCampaign: task_campaign.TaskCampaign{CampaignId: int64(campId)}}, false)

        assert.Equal(t, len(tests), count)

        assert.Equal(t, len(tests), len(campResults))
        assert.NoError(t, err)
        assert.NoError(t, mock.ExpectationsWereMet())
    })
}

func TestRepository_FindBank(t *testing.T) {
    conn, mock, red, err := config.InitMockConn(true)
    if err != nil {
        t.Fatal(err)
    }
    now := time.Now()
    repo := task_campaign.NewRepository(conn, red)
    campaignId := 1
    query := "SELECT campaign_banks.campaign_id, campaign_banks.bank_code, min(campaign_banks.from_date) as from_date, max(campaign_banks.to_date) as to_date, min(campaign_vouchers.from_date) as camp_from_date, max(campaign_vouchers.to_date) as camp_to_date FROM campaign_banks JOIN `campaign_vouchers` ON campaign_banks.campaign_id = campaign_vouchers.campaign_id WHERE (campaign_banks.campaign_id = (1) AND campaign_banks.deleted_at IS NULL AND (campaign_banks.to_date IS NULL OR campaign_banks.to_date >= now())) GROUP BY campaign_banks.campaign_id, campaign_banks.bank_code"
    tests := []config.TestStruct{
        {
            Obj: task_campaign.CampaignBank{
                TaskCampaign: task_campaign.TaskCampaign{CampaignId: int64(campaignId),
                    DateRange: task_campaign.DateRange{
                        FromDate: &now,
                        ToDate:   nil,
                    },
                },
                BankCode: "VCB1",
            },
            Query:       query,
            EmptyError:  true,
            EmptyOutput: false,
        },
        {
            Obj: task_campaign.CampaignBank{
                TaskCampaign: task_campaign.TaskCampaign{CampaignId: int64(campaignId),
                    DateRange: task_campaign.DateRange{
                        FromDate: &now,
                        ToDate:   nil,
                    },
                },
                BankCode: "VCB2",
            },
            Query:       query,
            EmptyError:  true,
            EmptyOutput: false,
        },
    }
    t.Run(fmt.Sprintf("Test Select TaskCampaign"), func(t *testing.T) {
        if err != nil {
            t.Fatalf("Error  %v not expect when mock db", err)
        }

        item1 := tests[0].Obj.(task_campaign.CampaignBank)
        item2 := tests[1].Obj.(task_campaign.CampaignBank)
        var rows = sqlmock.NewRows([]string{"campaign_id", "bank_code", "from_date", "to_date"}).
            AddRow(item1.TaskCampaign.CampaignId, item1.BankCode, item1.FromDate, item1.ToDate).
            AddRow(item2.TaskCampaign.CampaignId, item2.BankCode, item2.FromDate, item2.ToDate)

        mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
        var count, campResults, err = repo.GetTaskCampaigns(task_campaign.CampaignBank{TaskCampaign: task_campaign.TaskCampaign{CampaignId: int64(campaignId)}}, true)

        assert.Equal(t, len(tests), count)
        assert.Equal(t, len(tests), len(campResults))
        assert.NoError(t, err)
        assert.NoError(t, mock.ExpectationsWereMet())
    })
}
*/
