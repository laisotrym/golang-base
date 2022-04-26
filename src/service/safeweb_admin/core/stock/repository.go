//go:generate mockery --dir . --name IRepository --output ../../tests/service/user/mocks --outpkg userMocks --structname Repository --filename repository.go
package stock

import (
	"github.com/gocraft/dbr/v2"

	"safeweb.app/model"
)

type (
	IRepository interface {
		Get(string) ([]*model.Stock, error)
	}

	Repository struct {
		conn *dbr.Connection
	}
)

func NewRepository(conn *dbr.Connection) *Repository {
	return &Repository{conn: conn}
}

func (r *Repository) Get(name string) ([]*model.Stock, error) {
	var db model.Stock
	var stock []*model.Stock

	_, err := r.conn.NewSession(nil).Select(db.GetColumns()...).From(db.TableName()).
		Where("status <> ?", "done").
		Where("notify <> ?", "nope").
		Load(&stock)
	if err != nil {
		return nil, err
	} else {
		return stock, nil
	}
}
