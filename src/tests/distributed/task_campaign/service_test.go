package task_campaign_test

/*import (
    "context"
    "database/sql"
    "fmt"
    "net"
    "path"
    "testing"
    "time"

    "github.com/gocraft/dbr/v2"
    "github.com/pkg/errors"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    taskCampaignMocks "safeweb.app/tests/distributed/task_campaign/mocks"

    "safeweb.app/constant"
    "safeweb.app/distributed/task_campaign"
    "safeweb.app/model"
    "safeweb.app/pkg/helper"
    promo_error "safeweb.app/service/safeweb_lib/error"
)

func init() {
    err := promo_error.Init(path.Join(helper.RootDir(), "..", "..", "..", "..", "status.yml"))
    if err != nil {
        fmt.Print(err)
    }
}

func TestService_FlattenMerchant(t *testing.T) {
    ctx := context.Background()

    now := time.Now()
    mockCampObjs := []task_campaign.ITaskCampaign{
        task_campaign.CampaignMerchant{
            TaskCampaign: task_campaign.TaskCampaign{
                CampaignId: 1,
                DateRange: task_campaign.DateRange{
                    FromDate: &now,
                    ToDate:   &sql.NullTime{Time: now.Add(time.Duration(24 * 60 * 60 * time.Second)), Valid: true},
                },
            },
            MerchantCode: "0107305822",
        },
        task_campaign.CampaignMerchant{
            TaskCampaign: task_campaign.TaskCampaign{
                CampaignId: 1,
                DateRange: task_campaign.DateRange{
                    FromDate: &now,
                    ToDate:   &sql.NullTime{Time: now.Add(time.Duration(4 * 60 * 60 * time.Second)), Valid: true},
                },
            },
            MerchantCode: "0100692876",
        },
    }

    t.Run("db-fail-error", func(t *testing.T) {
        mockRepo := new(taskCampaignMocks.Repository)
        mockRepo.On("StoreRedisZRemAll", ctx, mock.Anything).
            Return(nil).Once()
        mockRepo.On("GetTaskCampaigns", mock.Anything, true).
            Return(0, nil, errors.New("Not found data")).Once()

        service := task_campaign.NewService(mockRepo)
        err := service.FlattenMerchant(ctx, 1)

        assert.Error(t, err)
        assert.EqualError(t, err, "Not found data")
    })

    t.Run("db-not-found", func(t *testing.T) {
        mockRepo := new(taskCampaignMocks.Repository)
        mockRepo.On("StoreRedisZRemAll", ctx, mock.Anything).
            Return(nil).Once()
        mockRepo.On("GetTaskCampaigns", mock.Anything, true).
            Return(0, nil, dbr.ErrNotFound).Once()

        service := task_campaign.NewService(mockRepo)
        err := service.FlattenMerchant(ctx, 1)

        assert.Error(t, err)
        assert.EqualError(t, err, dbr.ErrNotFound.Error())
    })

    t.Run("redis-fail-error", func(t *testing.T) {
        mockRepo := new(taskCampaignMocks.Repository)
        mockRepo.On("StoreRedisZRemAll", ctx, mock.Anything).
            Return(nil).Once()
        mockRepo.On("GetTaskCampaigns", mock.Anything, true).
            Return(2, mockCampObjs, nil).Once()

        mockRepo.On("StoreRedisHSet", ctx, mock.Anything, mock.Anything).
            Return(errors.New("redis: client is closed")).Once()

        mockRepo.On("StoreRedisZAdd", ctx, mock.Anything, mock.Anything, constant.RedisZAddTypeNX, false, mock.Anything).
            Return(errors.New("redis: client is closed")).Once()

        service := task_campaign.NewService(mockRepo)
        err := service.FlattenMerchant(ctx, 1)

        assert.Error(t, err)
        assert.EqualError(t, err, "redis: client is closed")
    })

    t.Run("redis-net-addr-error", func(t *testing.T) {
        mockRepo := new(taskCampaignMocks.Repository)
        mockRepo.On("StoreRedisZRemAll", ctx, mock.Anything).
            Return(nil).Once()
        mockRepo.On("GetTaskCampaigns", mock.Anything, true).
            Return(2, mockCampObjs, nil).Once()

        addrErr := net.AddrError{
            Err:  "Not Found",
            Addr: "127.0.0.1",
        }
        mockRepo.On("StoreRedisHSet", ctx, mock.Anything, mock.Anything).
            Return(&addrErr).Once()

        mockRepo.On("StoreRedisZAdd", ctx, mock.Anything, mock.Anything, constant.RedisZAddTypeNX, false, mock.Anything).
            Return(&addrErr).Once()

        service := task_campaign.NewService(mockRepo)
        err := service.FlattenMerchant(ctx, 1)

        assert.Error(t, err)
        assert.EqualError(t, err, addrErr.Error())
    })

    t.Run("success", func(t *testing.T) {
        mockRepo := new(taskCampaignMocks.Repository)
        mockRepo.On("StoreRedisZRemAll", ctx, mock.Anything).
            Return(nil).Once()
        mockRepo.On("GetTaskCampaigns", mock.Anything, true).
            Return(2, mockCampObjs, nil).Once()

        mockRepo.On("StoreRedisHSet", ctx, mock.Anything, mock.Anything).
            Return(nil).Once()

        mockRepo.On("StoreRedisZAdd", ctx, mock.Anything, mock.Anything, constant.RedisZAddTypeNX, false, mock.Anything).
            Return(nil).Once()

        // mockRepo.On("SetRedisExpireAt", ctx, mock.Anything, mock.AnythingOfType("time.Time")).
        //     Return(nil).Once()

        service := task_campaign.NewService(mockRepo)
        err := service.FlattenMerchant(ctx, 1)

        assert.NoError(t, err)
    })
}

func TestService_FlattenBaseInfo(t *testing.T) {
    ctx := context.Background()

    now := time.Now()
    mockCampObjs := []task_campaign.ITaskCampaign{
        task_campaign.CampaignVouchers{
            TaskCampaign: task_campaign.TaskCampaign{
                CampaignId: 1,
                DateRange: task_campaign.DateRange{
                    FromDate: &now,
                    ToDate:   &sql.NullTime{Time: now},
                },
            },
            VoucherCode: "VCB_001",
        },
        task_campaign.CampaignVouchers{
            TaskCampaign: task_campaign.TaskCampaign{
                CampaignId: 1,
                DateRange: task_campaign.DateRange{
                    FromDate: &now,
                    ToDate:   &sql.NullTime{Time: now},
                },
            },
            VoucherCode: "VCB_002",
        },
    }

    t.Run("db-fail-error-get-campaign", func(t *testing.T) {
        mockRepo := new(taskCampaignMocks.Repository)
        mockRepo.On("GetCampaign", mock.Anything).Return(nil, errors.New("error get")).Once()
        service := task_campaign.NewService(mockRepo)
        err := service.FlattenBaseInfo(ctx, 1)

        assert.Error(t, err)
        assert.EqualError(t, err, "Not found campaign: error get")
    })

    t.Run("db-fail-error-get-campaign-detail", func(t *testing.T) {
        mockRepo := new(taskCampaignMocks.Repository)
        mockRepo.On("GetCampaign", mock.Anything).Return(&model.Campaign{
            ID:   1,
            Name: "Campaign 01",
        }, nil).Once()
        mockRepo.On("GetCampaignDetail", mock.Anything).Return(nil, errors.New("error get campaign detail")).Once()
        service := task_campaign.NewService(mockRepo)

        err := service.FlattenBaseInfo(ctx, 1)

        assert.Error(t, err)
        assert.EqualError(t, err, "Not found Campaign Detail: error get campaign detail")
    })

    t.Run("db-fail-error-get-vouchers", func(t *testing.T) {
        // now := time.Now()
        mockRepo := new(taskCampaignMocks.Repository)

        mockRepo.On("GetCampaign", mock.Anything).Return(&model.Campaign{
            ID:   1,
            Name: "Campaign 01",
        }, nil).Once()
        mockRepo.On("GetCampaignDetail", mock.Anything).Return(&model.CampaignDetail{
            CampaignID:     1,
            QrTypeList:     "01, 02",
            PayMethodList:  "01, 02",
            TypeSourceList: "01, 02",
        }, nil).Once()
        mockRepo.On("GetTaskCampaigns", mock.Anything, mock.Anything).Return(0, nil, errors.New("Not found Campaign Voucher")).Once()
        service := task_campaign.NewService(mockRepo)

        err := service.FlattenBaseInfo(ctx, 1)

        assert.Error(t, err)
        assert.EqualError(t, err, "Not found Campaign Voucher")
    })

    t.Run("fail-to-connect-redis", func(t *testing.T) {
        mockRepo := new(taskCampaignMocks.Repository)

        mockRepo.On("GetCampaign", mock.Anything).Return(&model.Campaign{
            ID:   1,
            Name: "Campaign 01",
        }, nil).Once()

        mockRepo.On("GetCampaignDetail", mock.Anything).Return(&model.CampaignDetail{
            CampaignID:     1,
            QrTypeList:     "01, 02",
            PayMethodList:  "01, 02",
            TypeSourceList: "01, 02",
        }, nil).Once()

        mockRepo.On("GetTaskCampaigns", mock.Anything, mock.Anything).Return(2, mockCampObjs, nil).Once()

        mockRepo.On("StoreRedisSet", ctx, mock.Anything, mock.Anything, int64(0)).
            Return(errors.New("Failed To connect redis")).Once()

        service := task_campaign.NewService(mockRepo)

        err := service.FlattenBaseInfo(ctx, 1)

        assert.Error(t, err)
        assert.EqualError(t, err, "Failed To connect redis")
    })

    t.Run("success", func(t *testing.T) {
        mockRepo := new(taskCampaignMocks.Repository)

        mockRepo.On("GetCampaign", mock.Anything).Return(&model.Campaign{
            ID:   1,
            Name: "Campaign 01",
        }, nil).Once()

        mockRepo.On("GetCampaignDetail", mock.Anything).Return(&model.CampaignDetail{
            CampaignID:     1,
            QrTypeList:     "02",
            PayMethodList:  "01, 02",
            TypeSourceList: "01, 02",
        }, nil).Once()

        mockRepo.On("GetTaskCampaigns", mock.Anything, mock.Anything).Return(2, mockCampObjs, nil).Once()

        mockRepo.On("StoreRedisSet", ctx, mock.Anything, mock.Anything, int64(0)).Return(nil).Twice()

        service := task_campaign.NewService(mockRepo)

        err := service.FlattenBaseInfo(ctx, 1)

        assert.NoError(t, err)
    })

}
*/
