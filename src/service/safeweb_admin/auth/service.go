//go:generate mockery --dir . --name IService --output ../../tests/service/auth/mocks --outpkg authMocks --structname Service --filename service.go
package auth

import (
	"context"
	"log"

	"github.com/gogo/status"
	"google.golang.org/grpc/codes"
)

type (
	IService interface {
		Login(ctx context.Context, username string, password string) (*ItemUser, string, error)
	}

	Service struct {
		rp      IRepository
		jm      *JWTManager
		service IService
	}
)

func NewService(rp IRepository, jm *JWTManager) *Service {
	return &Service{rp: rp, jm: jm}
}

func (s *Service) Login(ctx context.Context, username string, password string) (*ItemUser, string, error) {
	if user, err := s.rp.Find(ctx, username); err != nil {
		log.Println("[ERROR] Login - rp.Find:", err)
		return nil, "", status.Error(codes.Unauthenticated, "incorrect username/password")
	} else if err := user.ComparePassword(password); err != nil {
		log.Println("[ERROR] Login - user.ComparePassword:", err)
		return nil, "", status.Error(codes.Unauthenticated, "incorrect username/password")
	} else {
		role := "user"
		if user.IsSuper {
			role = "super"
		} else if user.IsAdmin {
			role = "admin"
		}

		if token, err := s.jm.Generate(user, role); err != nil {
			log.Println("[ERROR] Login - jm.Generate:", err)
			return nil, "", status.Error(codes.Canceled, "cannot generate access token")
		} else {
			itemUser := &ItemUser{
				Id:       user.Id,
				UserName: user.UserName,
				FullName: user.FullName,
				Email:    user.Email,
				Phone:    user.Phone,
				Country:  user.Country,
				TimeZone: user.TimeZone,
			}
			return itemUser, token, err
		}
	}
}
