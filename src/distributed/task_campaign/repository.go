//go:generate mockery --dir . --name IRepository --output ../../tests/distributed/task_campaign/mocks --outpkg taskCampaignMocks --structname Repository --filename repository.go
package task_campaign

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"

	constant "safeweb.app/constant"
	"safeweb.app/model"
	"safeweb.app/model/enum"

	"github.com/go-redis/redis/v8"
	"github.com/gocraft/dbr/v2"
)

type (
	IRepository interface {
		GetCampaign(ITaskCampaign) (*model.Campaign, error)
		GetCampaignDetail(ITaskCampaign) (*model.CampaignDetail, error)
		GetTaskCampaigns(ITaskCampaign, bool) (int, []ITaskCampaign, error)
		StoreRedisHSet(context.Context, ITaskCampaign, map[string]interface{}, int64) error
		StoreRedisSet(context.Context, ITaskCampaign, interface{}, int64) error
		SetRedisExpireAt(context.Context, string, time.Time) error
		StoreRedisZAdd(context.Context, ITaskCampaign, []ZData, constant.RedisZAddType, bool, int64) error
		StoreRedisZRemAll(context.Context, ITaskCampaign) error
	}

	Repository struct {
		conn  *dbr.Connection
		redis *redis.Client
	}
)

func (r *Repository) GetCampaign(obj ITaskCampaign) (*model.Campaign, error) {
	campaignObj := model.Campaign{}

	err := r.conn.NewSession(nil).
		Select(obj.Columns(obj)...).
		From(obj.TableName()).
		Where("id = ? and deleted_at IS NULL and status = ?", obj.GetCampaignId(), enum.Approved).
		LoadOne(&campaignObj)

	return &campaignObj, err
}

func (r *Repository) GetCampaignDetail(obj ITaskCampaign) (*model.CampaignDetail, error) {
	campaignDetailObj := model.CampaignDetail{}

	err := r.conn.NewSession(nil).
		Select(obj.Columns(obj)...).
		From(obj.TableName()).
		Where("campaign_id = ? and deleted_at IS NULL", obj.GetCampaignId()).
		LoadOne(&campaignDetailObj)

	return &campaignDetailObj, err
}

func (r *Repository) GetTaskCampaigns(obj ITaskCampaign, isJoinVoucher bool) (int, []ITaskCampaign, error) {
	sess := r.conn.NewSession(nil)
	objTable := obj.TableName()
	var sel *dbr.SelectStmt
	if isJoinVoucher {
		campVouchers := &CampaignVouchers{}
		voucherTable := campVouchers.TableName()
		joinOn := fmt.Sprintf("%s.campaign_id = %s.campaign_id", objTable, voucherTable)

		sel = sess.Select(obj.VoucherCodeColumns(obj)...).Join(voucherTable, joinOn)
	} else {
		sel = sess.Select(obj.Columns(obj)...)
	}
	sel = sel.From(objTable).
		GroupBy(obj.GroupColumns(obj)...).
		Where(obj.WhereConditions(obj), obj.WhereArgs())

	var m []ITaskCampaign
	count, err := sel.Load(dbr.InterfaceLoader(&m, obj))
	return count, m, err
}

func (r *Repository) StoreRedisHSet(ctx context.Context, obj ITaskCampaign, values map[string]interface{}, secExpire int64) (err error) {
	rKey := obj.KeyRedis()

	defer func() {
		if err == nil && secExpire > 0 {
			err = r.redis.ExpireAt(ctx, rKey, time.Unix(secExpire, 0)).Err()
		}
	}()

	err = r.redis.HSet(ctx, rKey, values).Err()

	return
}

func (r *Repository) StoreRedisSet(ctx context.Context, obj ITaskCampaign, value interface{}, expiration int64) error {
	rKey := obj.KeyRedis()

	b, err := json.Marshal(&value)
	if err != nil {
		return err
	}
	return r.redis.Set(ctx, rKey, string(b), time.Duration(expiration)).Err()
}

func (r *Repository) SetRedisExpireAt(ctx context.Context, key string, expireAt time.Time) error {
	return r.redis.ExpireAt(ctx, key, expireAt).Err()
}

func (r *Repository) StoreRedisZAdd(ctx context.Context, obj ITaskCampaign, values []ZData, typ constant.RedisZAddType, rTyp bool, secExpire int64) (err error) {
	rKey := obj.KeyRedis()

	defer func() {
		if err == nil && secExpire > 0 {
			err = r.redis.ExpireAt(ctx, rKey, time.Unix(secExpire, 0)).Err()
		}
	}()

	dataLen := len(values)
	if dataLen > 0 {
		data := make([]*redis.Z, dataLen)
		for i, v := range values {
			data[i] = &redis.Z{
				Score:  float64(v.Score),
				Member: v.Member,
			}
		}

		if typ == constant.RedisZAddTypeNX && rTyp {
			err = r.redis.ZAddNXCh(ctx, rKey, data...).Err()
		} else if typ == constant.RedisZAddTypeNX {
			err = r.redis.ZAddNX(ctx, rKey, data...).Err()
		} else if typ == constant.RedisZAddTypeXX && rTyp {
			err = r.redis.ZAddXXCh(ctx, rKey, data...).Err()
		} else if typ == constant.RedisZAddTypeXX {
			err = r.redis.ZAddXX(ctx, rKey, data...).Err()
		} else {
			err = r.redis.ZAdd(ctx, rKey, data...).Err()
		}
	} else {
		err = errors.New("Not found data for store")
	}

	return
}

func (r *Repository) StoreRedisZRemAll(ctx context.Context, obj ITaskCampaign) error {
	return r.redis.ZRemRangeByRank(ctx, obj.KeyRedis(), int64(0), int64(-1)).Err()
}

func NewRepository(conn *dbr.Connection, redis *redis.Client) *Repository {
	return &Repository{conn: conn, redis: redis}
}
