package auth

import (
	"fmt"
	"reflect"

	"github.com/gogo/status"
	"google.golang.org/grpc/codes"
	"gopkg.in/go-playground/validator.v9"
	"safeweb.app/rpc/safeweb_admin"
	helper "safeweb.app/service/safeweb_lib/validate"
)

type (
	User struct {
		Name     string `json:"userName" validate:"username"`
		Password string `json:"password" validate:"required"`
		// Token    string `json:"token"`
		Token string `json:"token" validate:"captcha"` // <-- a custom validation rule
	}
)

func (user *User) clone(req *safeweb_admin.LoginRequest) {
	user.Name = req.GetUsername()
	user.Password = req.GetPassword()
	user.Token = req.GetToken()
}

func (user *User) validate() error {
	v := validator.New()

	_ = v.RegisterValidation("captcha", helper.ValidateCaptcha)
	_ = v.RegisterValidation("username", helper.ValidateUserName)

	err := v.Struct(user)
	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return fmt.Errorf("Invalid Argument - %v", err)
		}

		reflected := reflect.ValueOf(user).Elem()

		for _, e := range err.(validator.ValidationErrors) {
			field, _ := reflected.Type().FieldByName(e.StructField())
			name := field.Tag.Get("json")

			switch name {
			case "userName":
				return status.Error(codes.InvalidArgument, "The user name must be at least 3 characters long and must not have any whitespaces.")
			case "email":
				return status.Error(codes.InvalidArgument, "The email is invalid.")
			case "token":
				return status.Error(codes.InvalidArgument, "The token is invalid.")
			default:
				return status.Errorf(codes.InvalidArgument, "Invalid Argument - %v", e)
			}
		}
	}
	return nil
}
