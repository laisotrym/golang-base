package user

import (
	"fmt"
	"reflect"

	"gopkg.in/go-playground/validator.v9"
	"safeweb.app/rpc/safeweb_admin"
	helper "safeweb.app/service/safeweb_lib/validate"
)

type (
	User struct {
		Name     string `json:"userName" validate:"username"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"strongpassword"` // <-- a custom validation rule
		FullName string `json:"fullName"`
		Country  string `json:"country"`
		TimeZone string `json:"timeZone"`
		Phone    string `json:"phone"`
		Token    string `json:"token" validate:"captcha"`
	}
)

func (user *User) clone(req *safeweb_admin.SignUpReq) {
	user.Name = req.GetUserName()
	user.Email = req.GetEmail()
	user.Password = req.GetPassword()
	user.FullName = req.GetFullName()
	user.Country = req.GetCountry()
	user.TimeZone = req.GetTimeZone()
	user.Phone = req.GetPhone()
	user.Token = req.GetToken()
}

func (user *User) validate() error {
	v := validator.New()

	_ = v.RegisterValidation("strongpassword", helper.ValidateStrongPassword)
	_ = v.RegisterValidation("captcha", helper.ValidateCaptcha)
	_ = v.RegisterValidation("phone", helper.ValidatePhone)
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
				return fmt.Errorf("The user name must be at least 3 characters long and must not have any whitespaces.")
			case "email":
				return fmt.Errorf("The email is invalid.")
			case "password":
				return fmt.Errorf("The password must have at least 8 characters, including one number, one lowercase character, one uppercase character, one special character.")
			case "phone":
				return fmt.Errorf("The phone number is invalid.")
			case "token":
				return fmt.Errorf("The token is invalid.")
			default:
				return fmt.Errorf("Invalid Argument - %v", e)
			}
		}
	}
	return nil
}
