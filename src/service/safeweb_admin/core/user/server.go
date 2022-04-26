package user

import (
	"context"
	"database/sql"

	"github.com/gogo/protobuf/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"safeweb.app/rpc/safeweb_admin"
	safeweb_lib_helper "safeweb.app/service/safeweb_lib/helper"
)

type (
	// UserServer is the server for authentication
	Server struct {
		safeweb_admin.UnimplementedUserServiceServer
		service IService
	}
)

// NewAuthServer returns a new auth server
func NewServer(service IService) *Server {
	return &Server{service: service}
}

func (s *Server) SignUp(ctx context.Context, req *safeweb_admin.SignUpReq) (*safeweb_admin.SignUpRes, error) {
	res := &safeweb_admin.SignUpRes{}

	// Validate inputs
	user := User{}
	user.clone(req)
	err := user.validate()
	if err != nil {
		res.Success = false
		res.Message = err.Error()
		return res, nil
	}

	// Decrypt password
	realPass, err := safeweb_lib_helper.Decrypt(user.Password)
	if err != nil {
		res.Success = false
		res.Message = err.Error()
		return res, nil
	}

	// log.Println("Password decrypted:", realPass)
	user.Password = realPass

	// Create an account
	if _, err := s.service.Create(ctx, &ItemUser{
		UserName: user.Name,
		Password: user.Password,
		Email:    user.Email,
		Phone:    sql.NullString{String: user.Phone, Valid: true},
		FullName: sql.NullString{String: user.FullName, Valid: true},
		Country:  sql.NullString{String: user.Country, Valid: true},
		TimeZone: sql.NullString{String: user.TimeZone, Valid: true},
	}); err != nil {
		res.Success = false
		res.Message = "Create an account fail. This error could be caused by an existing username, email or phone number."
		return res, nil
	}

	res.Success = true
	res.Message = "Success"

	return res, nil
}

func (s *Server) Activate(ctx context.Context, req *safeweb_admin.ActivateUserReq) (*safeweb_admin.ActivateUserRes, error) {
	res := &safeweb_admin.ActivateUserRes{}

	success, err := s.service.Confirm(ctx, req.GetToken())

	if err != nil {
		return nil, err
	}

	if success {
		res.Success = true
		res.Message = "Success"

		return res, nil
	}

	res.Success = false
	res.Message = "Invalid Request"
	return res, nil
}

func (s *Server) TestName(ctx context.Context, username *types.StringValue) (*safeweb_admin.TestNameRes, error) {
	res := &safeweb_admin.TestNameRes{}
	exist, _ := s.service.ExistName(ctx, username.Value)

	if exist {
		res.Success = true
		res.Message = "Name Exists"

		return res, nil
	}

	res.Success = false
	res.Message = "Not Found"
	return res, nil
}

func (s *Server) Add(ctx context.Context, req *safeweb_admin.AddUserReq) (*safeweb_admin.User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}

func (s *Server) Get(ctx context.Context, id *types.Int64Value) (*safeweb_admin.User, error) {
	//validate input
	userRes := &safeweb_admin.User{}

	user, err := s.service.Get(ctx, id.GetValue())
	if err != nil {
		return nil, err
	}

	if user != nil {
		userRes.Id = user.Id
		userRes.UserName = user.UserName
		userRes.Email = user.Email
		userRes.Phone = user.Phone.String
		userRes.FullName = user.FullName.String
		userRes.Status = user.Status
		userRes.IsSuper = user.IsSuper
		userRes.IsAdmin = user.IsAdmin
		userRes.CreatedBy = user.CreatedBy
		userRes.CreatedAt = user.CreatedAt
		userRes.UpdatedBy = user.UpdatedBy
		userRes.UpdatedAt = user.UpdatedAt
	}

	return userRes, nil
}

func (s *Server) ForgotPassword(ctx context.Context, req *safeweb_admin.ForgotPasswordReq) (*safeweb_admin.ForgotPasswordRes, error) {
	res := &safeweb_admin.ForgotPasswordRes{}
	res.Success = false
	res.Message = "Not Found Or Not Member"

	// TODO: Validate req
	ok, err := s.service.ResetPassword(ctx, req.GetEmail())
	if err != nil {
		return nil, err
	}
	if ok {
		res.Success = true
		res.Message = "Success"
	}

	return res, nil
}

func (s *Server) NewPassword(ctx context.Context, req *safeweb_admin.NewPasswordReq) (*safeweb_admin.NewPasswordRes, error) {
	res := &safeweb_admin.NewPasswordRes{}

	// Decrypt password
	realPass, err := safeweb_lib_helper.Decrypt(req.GetPassword())
	if err != nil {
		return nil, err
	}

	// log.Println("Password decrypted:", realPass)

	success, err := s.service.CreatePassword(ctx, realPass, req.GetToken())

	if err != nil {
		return nil, err
	}

	if success {
		res.Success = true
		res.Message = "Success"

		return res, nil
	}

	res.Success = false
	res.Message = "Invalid Request"
	return res, nil
}

func (s *Server) Update(ctx context.Context, req *safeweb_admin.UpdateUserReq) (*safeweb_admin.User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (s *Server) Delete(ctx context.Context, id *types.Int64Value) (*safeweb_admin.DeleteUserRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (s *Server) List(ctx context.Context, req *safeweb_admin.ListUserReq) (*safeweb_admin.ListUserRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HOho not implemented")
}
