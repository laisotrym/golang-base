//go:generate mockery --dir . --name IRepository --output ../../tests/service/auth/mocks --outpkg authMocks --structname Repository --filename repository.go
package auth

import (
	"context"

	"github.com/gocraft/dbr/v2"

	"safeweb.app/model"
)

type (
	IRepository interface {
		// Save saves a user to the store
		Save(ctx context.Context, user *model.User) (*model.User, error)
		// Find finds a user by username
		Find(ctx context.Context, username string) (*model.User, error)
	}

	Repository struct {
		conn *dbr.Connection
	}
)

func NewRepository(conn *dbr.Connection) *Repository {
	return &Repository{conn: conn}
}

// Save saves a user to the store
func (r *Repository) Save(ctx context.Context, user *model.User) (*model.User, error) {
	_, err := r.conn.NewSession(nil).
		InsertInto(user.TableName()).
		Columns(user.T().InsertColumns()...).
		Record(user).
		Exec()

	return user, err
}

// Find finds a user by username
func (r *Repository) Find(ctx context.Context, username string) (*model.User, error) {
	sess := r.conn.NewSession(nil)
	user := model.User{}
	err := sess.Select("id", "username", "email", "phone", "full_name", "password_hash", "country", "time_zone", "status", "is_super", "is_admin").From(user.TableName()).Where("username = ?", username).Where("is_member = ?", true).LoadOne(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
