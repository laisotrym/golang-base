//go:generate mockery --dir . --name IRepository --output ../../tests/service/user/mocks --outpkg userMocks --structname Repository --filename repository.go
package user

import (
	"github.com/gocraft/dbr/v2"

	"safeweb.app/.third_party/github.com/fatih/structs"
	"safeweb.app/model"
	"safeweb.app/model/enum"
)

type (
	IRepository interface {
		Save(*model.User) (int64, error)
		Get(int64) (*model.User, error)
		Update(*model.User) (*int64, error)
		Delete(int64) (*int64, error)
		List(*SearchInput) ([]*model.User, error)
		CountList(*SearchInput) (int64, error)
		All() ([]*model.User, error)
		CountAll() (int64, error)
		GetByUsername(string) (*model.User, error)
		SetMember(string) (*int64, error)
		SetForgotPassword(string, string, string) (*int64, error)
		UpdatePassword(string, string) (*int64, error)
	}

	Repository struct {
		conn *dbr.Connection
	}
)

func NewRepository(conn *dbr.Connection) *Repository {
	return &Repository{conn: conn}
}

func (r *Repository) Save(user *model.User) (int64, error) {
	if result, err := r.conn.NewSession(nil).
		InsertInto(user.TableName()).
		Columns(user.T().InsertColumns()...).
		Record(user).
		Exec(); err != nil {
		return 0, err
	} else {
		return result.LastInsertId()
	}
}

func (r *Repository) Get(id int64) (*model.User, error) {
	var user model.User
	sess := r.conn.NewSession(nil)
	if err := sess.Select(user.GetColumns()...).From(user.TableName()).Where("id = ?", id).LoadOne(&user); err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

func (r *Repository) Update(user *model.User) (*int64, error) {
	if result, err := r.conn.NewSession(nil).
		Update(user.TableName()).
		SetMap(structs.Map(user)).
		Where("id = ?", user.Id).
		Exec(); err != nil {
		return nil, err
	} else {
		count, err := result.RowsAffected()
		return &count, err
	}
}

func (r *Repository) Delete(id int64) (*int64, error) {
	var user model.User
	if result, err := r.conn.NewSession(nil).
		DeleteFrom(user.TableName()).
		Where("id = ?", id).
		Exec(); err != nil {
		return nil, err
	} else {
		count, err := result.RowsAffected()
		return &count, err
	}
}

func (r *Repository) buildListQuery(stmt *dbr.SelectStmt, cond *SearchInput) {
	stmt.Where("username LIKE ?", "%"+cond.UserName+"%").
		Where("email LIKE ?", "%"+cond.Email+"%").
		Where("phone LIKE ?", "%"+cond.Phone+"%").
		Where("full_name LIKE ?", "%"+cond.FullName+"%")

	if cond.Id > 0 {
		stmt.Where("id = ?", cond.Id)
	}
}

func (r *Repository) List(cond *SearchInput) ([]*model.User, error) {
	var user model.User
	var users []*model.User
	sess := r.conn.NewSession(nil)
	selectStmt := sess.Select(user.GetColumns()...).From(user.TableName())

	r.buildListQuery(selectStmt, cond)

	if _, err := selectStmt.Paginate(uint64(cond.Page+1), uint64(cond.PageSize)).OrderAsc("id").Load(&users); err != nil {
		return nil, err
	} else {
		return users, nil
	}
}

func (r *Repository) CountList(cond *SearchInput) (int64, error) {
	var user model.User
	sess := r.conn.NewSession(nil)
	countStmt := sess.Select("COUNT(1)").From(user.TableName())

	r.buildListQuery(countStmt, cond)

	if total, err := countStmt.ReturnInt64(); err != nil {
		return 0, err
	} else {
		return total, nil
	}
}

func (r *Repository) All() ([]*model.User, error) {
	var user model.User
	var users []*model.User
	_, err := r.conn.NewSession(nil).Select(user.GetColumns()...).From(user.TableName()).OrderAsc("id").Load(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) CountAll() (int64, error) {
	var user model.User
	if total, err := r.conn.NewSession(nil).Select("COUNT(1)").From(user.TableName()).ReturnInt64(); err != nil {
		return 0, err
	} else {
		return total, nil
	}
}

func (r *Repository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	sess := r.conn.NewSession(nil)
	err := sess.Select(user.GetColumns()...).From(user.TableName()).Where("username = ?", username).LoadOne(&user)
	if err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

func (r *Repository) SetMember(token string) (*int64, error) {
	if result, err := r.conn.NewSession(nil).
		Update("core_user").
		Set("is_member", true).
		Set("notif", nil).
		Set("updated_by", enum.UpdatedByServer).
		Where("confirm_token = ?", token).
		Exec(); err != nil {
		return nil, err
	} else {
		count, err := result.RowsAffected()
		return &count, err
	}
}

func (r *Repository) SetForgotPassword(token string, time string, email string) (*int64, error) {
	if result, err := r.conn.NewSession(nil).
		Update("core_user").
		Set("notif", enum.ResetPassword).
		Set("confirm_token", token).
		Set("confirm_time", time).
		Set("updated_by", enum.UpdatedByAdmin).
		Where("email = ?", email).
		Where("is_member = ?", true).
		Exec(); err != nil {
		return nil, err
	} else {
		count, err := result.RowsAffected()
		return &count, err
	}
}

func (r *Repository) UpdatePassword(passwordHash string, token string) (*int64, error) {
	if result, err := r.conn.NewSession(nil).
		Update("core_user").
		Set("password_hash", passwordHash).
		Set("notif", nil).
		Set("updated_by", enum.UpdatedByServer).
		Where("confirm_token = ?", token).
		Exec(); err != nil {
		return nil, err
	} else {
		count, err := result.RowsAffected()
		return &count, err
	}
}
