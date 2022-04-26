package stock

import (
	"context"

	"github.com/gogo/status"
	"google.golang.org/grpc/codes"
	"safeweb.app/rpc/safeweb_admin"
)

type (
	Server struct {
		safeweb_admin.UnimplementedStockServiceServer
		service IService
	}
)

func NewServer(service IService) *Server {
	return &Server{service: service}
}

func (s *Server) GetAll(ctx context.Context, req *safeweb_admin.GetAllReq) (*safeweb_admin.GetAllRes, error) {
	res := &safeweb_admin.GetAllRes{}

	// Validate inputs
	user := User{}
	user.clone(req)
	err := user.validate()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	stockList, err := s.service.Get(ctx, user.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if stockList != nil {
		var stocks []*safeweb_admin.Stock
		count := len(stockList)
		for i := 0; i < count; i++ {
			stock := &safeweb_admin.Stock{
				Id:         stockList[i].Id,
				MaCk:       stockList[i].MaCk,
				GiaMua:     stockList[i].GiaMua,
				KhoiLuong:  stockList[i].KhoiLuong,
				SoNgay:     stockList[i].SoNgay,
				TienLo:     stockList[i].TienLo,
				TyLeLo:     stockList[i].TyLeLo,
				GiaMax:     stockList[i].GiaMax,
				LaiMax:     stockList[i].LaiMax,
				TyLeMax:    stockList[i].TyLeMax,
				GiaHomNay:  stockList[i].GiaHomNay,
				LaiHomNay:  stockList[i].LaiHomNay,
				TyLeHomNay: stockList[i].TyLeHomNay,
				TrangThai:  stockList[i].TrangThai,
			}
			stocks = append(stocks, stock)
		}

		res.StockList = stocks
		res.Total = int64(count)
		return res, nil
	}

	return nil, nil
}
