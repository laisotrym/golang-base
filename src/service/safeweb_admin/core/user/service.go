//go:generate mockery --dir . --name IService --output ../../tests/service/user/mocks --outpkg userMocks --structname Service --filename service.go
package user

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"safeweb.app/model"
	"safeweb.app/model/enum"
	"safeweb.app/rpc/safeweb_lib"
	safeweb_lib_utils "safeweb.app/service/safeweb_lib/utils"
)

type (
	IService interface {
		Create(context.Context, *ItemUser) (*ItemUser, error)
		Get(context.Context, int64) (*ItemUser, error)
		Update(context.Context, *ItemUser) (*ItemUser, error)
		Delete(context.Context, int64) error
		List(context.Context, *SearchInput) (*ListOutput, error)
		ListAll(context.Context) (*ListOutput, error)
		Confirm(context.Context, string) (bool, error)
		ExistName(context.Context, string) (bool, error)
		ResetPassword(context.Context, string) (bool, error)
		CreatePassword(context.Context, string, string) (bool, error)
	}

	Service struct {
		rp      IRepository
		service IService
	}
)

func NewService(rp IRepository) *Service {
	return &Service{rp: rp}
}

func (s Service) Create(ctx context.Context, user *ItemUser) (*ItemUser, error) {
	if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost); err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	} else {
		loc, _ := time.LoadLocation("UTC")
		timeStr := time.Now().In(loc).Format(enum.TimeZoneLayout)
		token := safeweb_lib_utils.NewShortID()

		modelUser := &model.User{
			UserName:     user.UserName,
			PasswordHash: string(hashedPassword),
			Email:        user.Email,
			Phone:        user.Phone,
			FullName:     user.FullName,
			Country:      user.Country,
			TimeZone:     user.TimeZone,
			Status:       safeweb_lib.ActiveStatus_INACTIVE,
			IsSuper:      false,
			IsAdmin:      true,
			IsMember:     false,
			Notif:        enum.SignUp,
			ConfirmToken: sql.NullString{String: token, Valid: true},
			ConfirmTime:  sql.NullString{String: timeStr, Valid: true},
		}
		if userId, err := s.rp.Save(modelUser); err != nil {
			log.Println("[ERROR] Save new account to database - user name:", user.UserName, "email:", user.Email, "phone:", user.Phone, "err:", err)
			return nil, err
		} else {
			t1, _ := time.Parse(enum.TimeZoneLayout, modelUser.ConfirmTime.String)

			return &ItemUser{
				Id:           userId,
				UserName:     modelUser.UserName,
				Email:        modelUser.Email,
				Phone:        modelUser.Phone,
				FullName:     modelUser.FullName,
				Country:      modelUser.Country,
				TimeZone:     modelUser.TimeZone,
				Status:       modelUser.Status,
				IsSuper:      modelUser.IsSuper,
				IsAdmin:      modelUser.IsAdmin,
				IsMember:     modelUser.IsMember,
				Notif:        sql.NullString{String: modelUser.Notif, Valid: true},
				ConfirmToken: modelUser.ConfirmToken.String,
				ConfirmTime:  t1,
				CreatedBy:    modelUser.CreatedBy.Int64,
				CreatedAt:    modelUser.CreatedAt.Time.Unix(),
				UpdatedBy:    modelUser.UpdatedBy.Int64,
				UpdatedAt:    modelUser.UpdatedAt.Time.Unix(),
			}, nil
		}
	}
}

func (s Service) Get(ctx context.Context, i int64) (*ItemUser, error) {
	if modelUser, err := s.rp.Get(i); err != nil && err != dbr.ErrNotFound {
		return nil, errors.Wrap(err, "error with code")
	} else if modelUser == nil {
		return nil, status.Error(codes.NotFound, "not found user")
	} else {
		return &ItemUser{
			Id:        modelUser.Id,
			UserName:  modelUser.UserName,
			Email:     modelUser.Email,
			Phone:     modelUser.Phone,
			FullName:  modelUser.FullName,
			Status:    modelUser.Status,
			IsSuper:   modelUser.IsSuper,
			IsAdmin:   modelUser.IsAdmin,
			CreatedBy: modelUser.CreatedBy.Int64,
			CreatedAt: modelUser.CreatedAt.Time.Unix(),
			UpdatedBy: modelUser.UpdatedBy.Int64,
			UpdatedAt: modelUser.UpdatedAt.Time.Unix(),
		}, nil
	}
}

func (s Service) Update(ctx context.Context, user *ItemUser) (*ItemUser, error) {
	panic("implement me")
}

func (s Service) Delete(ctx context.Context, i int64) error {
	panic("implement me")
}

func (s Service) List(ctx context.Context, input *SearchInput) (*ListOutput, error) {
	panic("implement me")
}

func (s Service) ListAll(ctx context.Context) (*ListOutput, error) {
	panic("implement me")
}

func (s Service) Confirm(ctx context.Context, token string) (bool, error) {

	count, err := s.rp.SetMember(token)

	if err != nil {
		log.Println("[ERROR] Activate account - token:", token, "err:", err)
		return false, err
	}

	if *count == 0 {
		return false, nil
	}

	return true, nil
}

func (s Service) ExistName(ctx context.Context, username string) (bool, error) {

	_, err := s.rp.GetByUsername(username)
	if err != nil {
		log.Println("[ERROR] Test account - user name:", username, "err:", err)
		return false, err
	}

	return true, nil
}

func (s Service) ResetPassword(ctx context.Context, email string) (bool, error) {
	loc, _ := time.LoadLocation("UTC")
	timeStr := time.Now().In(loc).Format(enum.TimeZoneLayout)
	token := safeweb_lib_utils.NewShortID()

	count, err := s.rp.SetForgotPassword(token, timeStr, email)
	if err != nil {
		log.Println("[ERROR] Reset password - email:", email, "err:", err)
		return false, err
	}

	if *count == 0 {
		return false, nil
	}

	return true, nil
}

func (s Service) CreatePassword(ctx context.Context, password string, token string) (bool, error) {
	if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err != nil {
		return false, fmt.Errorf("cannot hash password: %w", err)
	} else {
		count, err := s.rp.UpdatePassword(string(hashedPassword), token)
		if err != nil {
			log.Println("[ERROR] New password - token:", token, "err:", err)
			return false, err
		}

		if *count == 0 {
			return false, nil
		}

		return true, nil
	}
}
