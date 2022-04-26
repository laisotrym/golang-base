//go:generate wire gen
// +build wireinject

package user

import (
	"github.com/gocraft/dbr/v2"
	"github.com/google/wire"
)

func InitializeServer(conn *dbr.Connection) *Server {
	wire.Build(NewServer,
		NewService,
		NewRepository,
		wire.Bind(new(IService), new(*Service)),
		wire.Bind(new(IRepository), new(*Repository)))
	return &Server{}
}

func InitializeService(conn *dbr.Connection) *Service {
	wire.Build(NewService,
		NewRepository,
		wire.Bind(new(IRepository), new(*Repository)))
	return &Service{}
}
