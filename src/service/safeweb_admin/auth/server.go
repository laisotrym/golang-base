package auth

import (
	"context"

	"safeweb.app/rpc/safeweb_admin"
	safeweb_lib_helper "safeweb.app/service/safeweb_lib/helper"
)

type (
	// AuthServer is the server for authentication
	Server struct {
		safeweb_admin.UnimplementedAuthServiceServer
		service IService
	}
)

// NewAuthServer returns a new auth server
func NewServer(service IService) *Server {
	return &Server{service: service}
}

// Login is a unary RPC to login user
func (s Server) Login(ctx context.Context, req *safeweb_admin.LoginRequest) (*safeweb_admin.LoginResponse, error) {
	// Validate inputs
	inputs := User{}
	inputs.clone(req)
	if err := inputs.validate(); err != nil {
		return nil, err
	}

	// Decrypt password
	realPass, err := safeweb_lib_helper.Decrypt(inputs.Password)
	if err != nil {
		return nil, err
	}

	// log.Println("Password decrypted:", realPass)
	inputs.Password = realPass

	user, token, err := s.service.Login(ctx, inputs.Name, inputs.Password)

	if err == nil && user != nil {
		res := &safeweb_admin.LoginResponse{
			UserName:    user.UserName,
			FullName:    user.FullName.String,
			Email:       user.Email,
			Phone:       user.Phone.String,
			Country:     user.Country.String,
			Timezone:    user.TimeZone.String,
			AccessToken: token,
			FreshToken:  token,
		}

		return res, nil
	}

	return nil, err
}
