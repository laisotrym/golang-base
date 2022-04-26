//go:generate mockery --dir . --name IService --output ../../tests/service/user/mocks --outpkg userMocks --structname Service --filename service.go
package stock

import (
	"context"
	"log"

	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	IService interface {
		Get(context.Context, string) ([]*ItemStock, error)
	}

	Service struct {
		rp      IRepository
		service IService
	}
)

func NewService(rp IRepository) *Service {
	return &Service{rp: rp}
}

func (s Service) Get(ctx context.Context, user string) ([]*ItemStock, error) {
	result, err := s.rp.Get(user)
	if err != nil && err != dbr.ErrNotFound {
		log.Println("[ERROR]", err.Error())
		return nil, errors.Wrap(err, "error with code")
	}

	if result == nil {
		log.Println("[ERROR]", "not found data")
		return nil, status.Error(codes.NotFound, "not found data")
	}

	var stock *ItemStock
	var stocks []*ItemStock
	count := len(result)

	for i := 0; i < count; i++ {
		stock = &ItemStock{
			Id:         result[i].Id,
			MaCk:       result[i].MaCk,
			GiaMua:     result[i].GiaMua,
			KhoiLuong:  result[i].KhoiLuong,
			SoNgay:     result[i].SoNgay,
			TienLo:     result[i].TienLo,
			TyLeLo:     result[i].TyLeLo,
			GiaMax:     result[i].GiaMax,
			LaiMax:     result[i].LaiMax,
			TyLeMax:    result[i].TyLeMax,
			GiaHomNay:  result[i].GiaHomNay,
			LaiHomNay:  result[i].LaiHomNay,
			TyLeHomNay: result[i].TyLeHomNay,
			TrangThai:  result[i].TrangThai,
		}

		stocks = append(stocks, stock)
	}
	return stocks, nil
}
