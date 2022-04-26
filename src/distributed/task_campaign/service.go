//go:generate mockery --dir . --name IService --output ../../tests/distributed/task_campaign/mocks --outpkg taskCampaignMocks --structname Service --filename service.go
package task_campaign

import (
    "context"
    "reflect"
    "time"
    
    "github.com/pkg/errors"
    
    "safeweb.app/constant"
)

type IService interface {
    FlattenBaseInfo(context.Context, int64) error
    FlattenBank(context.Context, int64) error
    FlattenMerchant(context.Context, int64) error
    FlattenTerminal(context.Context, int64) error
    FlattenBenefit(context.Context, int64) error
}

type Service struct {
    rp IRepository
}

func (s Service) FlattenBaseInfo(ctx context.Context, campId int64) error {
    if campaign, err := s.rp.GetCampaign(Campaign{TaskCampaign: TaskCampaign{CampaignId: campId}}); err != nil {
        return errors.Wrap(err, "Not found campaign")
    } else if campaignDetail, err := s.rp.GetCampaignDetail(CampaignDetail{TaskCampaign: TaskCampaign{CampaignId: campId}}); err != nil {
        return errors.Wrap(err, "Not found Campaign Detail")
    } else if count, campResults, err := s.rp.GetTaskCampaigns(CampaignVouchers{TaskCampaign: TaskCampaign{CampaignId: campId}}, false); err != nil {
        return err
    } else if count <= 0 || campResults == nil || reflect.ValueOf(campResults).IsNil() {
        return errors.New("Not found Campaign Vouchers")
    } else {
        for _, obj := range campResults {
            var expireTime int64
            if obj.GetToDate() != nil && obj.GetToDate().Valid {
                expireTime = obj.GetToDate().Time.UnixNano() - time.Now().UnixNano()
            }
            
            if err = s.rp.StoreRedisSet(ctx, obj, &BaseInfo{
                CampaignId:     campaign.ID,
                CampaignName:   campaign.Name,
                VoucherCode:    *obj.FieldRedis(),
                FromDate:       obj.GetFromDate(),
                ToDate:         obj.GetToDate(),
                QrTypeList:     campaignDetail.GetQrTypeList(),
                TypeSourceList: campaignDetail.GetTypeSourceList(),
                PayMethodList:  campaignDetail.GetPaymentMethodList(),
            }, expireTime); err != nil {
                return err
            }
        }
    }
    return nil
}

func (s Service) FlattenBank(ctx context.Context, campId int64) error {
    return s.flattenProcess(ctx, CampaignBank{TaskCampaign: TaskCampaign{CampaignId: campId}}, safeweb_lib_constant.RedisTypeSortedSets, true)
}

func (s Service) FlattenMerchant(ctx context.Context, campId int64) error {
    return s.flattenProcess(ctx, CampaignMerchant{TaskCampaign: TaskCampaign{CampaignId: campId}}, safeweb_lib_constant.RedisTypeSortedSets, true)
}

func (s Service) FlattenTerminal(ctx context.Context, campId int64) error {
    return s.flattenProcess(ctx, CampaignTerminal{CampaignMerchant: CampaignMerchant{TaskCampaign: TaskCampaign{CampaignId: campId}}}, safeweb_lib_constant.RedisTypeSortedSets, true)
}

func (s Service) FlattenBenefit(ctx context.Context, campId int64) error {
    return s.flattenProcess(ctx, CampaignBenefit{TaskCampaign: TaskCampaign{CampaignId: campId}}, safeweb_lib_constant.RedisTypeSortedSets, true)
}

func (s Service) FlattenBlacklistCustomer(ctx context.Context, mobile string) error {
    return s.flattenProcess(ctx, BlackListCustomer{Mobile: mobile}, safeweb_lib_constant.RedisTypeSortedSets, false)
}

func (s Service) FlattenBlacklistTerminal(ctx context.Context, merchantCode string) error {
    return s.flattenProcess(ctx, BlackListTerminal{MerchantCode: merchantCode}, safeweb_lib_constant.RedisTypeSortedSets, false)
}

func (s Service) flattenProcess(ctx context.Context, objInfo ITaskCampaign, redTyp safeweb_lib_constant.RedisType, isJoinVoucher bool) error {
    if err := s.rp.StoreRedisZRemAll(ctx, objInfo); err != nil {
        return err
    } else if count, campResults, err := s.rp.GetTaskCampaigns(objInfo, isJoinVoucher); err != nil {
        return err
    } else if count <= 0 || campResults == nil || reflect.ValueOf(campResults).IsNil() {
        return errors.New("Not found data")
    } else {
        var dataSortedSets []ZData
        data := make(map[string]interface{}, 0)
        var expireAt int64
        
        for _, obj := range campResults {
            sortedScore := safeweb_lib_constant.MaxInt64 - int64(100)
            
            campToDate := obj.GetCampToDate()
            if campToDate != nil {
                sortedScore = campToDate.Unix()
                if reflect.ValueOf(expireAt).IsZero() || campToDate.Unix() > expireAt {
                    expireAt = campToDate.Unix()
                }
            }
            
            toDate := obj.GetToDate()
            if toDate != nil && toDate.Valid && toDate.Time.Unix() < sortedScore {
                sortedScore = toDate.Time.Unix()
            }
            
            if sortedScore > time.Now().Unix() {
                if memberObj, err := obj.DataRedisSortedSets(); err != nil {
                    return err
                } else if memberObj != nil {
                    score := obj.DataRedisSortedSetsScore(sortedScore)
                    dataSortedSets = append(dataSortedSets, ZData{
                        Score:  score,
                        Member: memberObj,
                    })
                }
                if field := obj.FieldRedis(); field != nil {
                    data[*field] = obj.DataRedis()
                }
            }
        }
        
        if len(data) > 0 || len(dataSortedSets) > 0 {
            if redTyp == safeweb_lib_constant.RedisTypeSortedSets && len(dataSortedSets) > 0 {
                return s.rp.StoreRedisZAdd(ctx, objInfo, dataSortedSets, safeweb_lib_constant.RedisZAddTypeNX, false, expireAt)
            } else if len(data) > 0 {
                return s.rp.StoreRedisHSet(ctx, objInfo, data, expireAt)
            }
        }
        return errors.New("Not found flatten data")
    }
}

func NewService(rp IRepository) *Service {
    return &Service{rp: rp}
}