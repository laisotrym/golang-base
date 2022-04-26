//go:generate wire gen
// +build wireinject

package task_campaign

import (
    "github.com/go-redis/redis/v8"
    "github.com/gocraft/dbr/v2"
    "github.com/google/wire"
)

func InitializeService(conn *dbr.Connection, redis *redis.Client) *Service {
    wire.Build(NewService,
        NewRepository,
        wire.Bind(new(IRepository), new(*Repository)))
    return &Service{}
}
