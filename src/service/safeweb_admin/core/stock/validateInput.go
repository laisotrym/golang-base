package stock

import (
	"fmt"
	"reflect"

	"gopkg.in/go-playground/validator.v9"
	"safeweb.app/rpc/safeweb_admin"
	helper "safeweb.app/service/safeweb_lib/validate"
)

type (
	User struct {
		Name  string `json:"username"`
		Token string `json:"token"`
	}
)

func (user *User) clone(req *safeweb_admin.GetAllReq) {
	user.Name = req.GetUsername()
	user.Token = req.GetToken()
}

func (user *User) validate() error {
	v := validator.New()

	_ = v.RegisterValidation("vlUsername", helper.ValidateUserName)

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
			case "username":
				return fmt.Errorf("The user name must be at least 3 characters long and must not have any whitespaces.")
			case "token":
				return fmt.Errorf("The token is invalid.")
			default:
				return fmt.Errorf("Invalid Argument - %v", e)
			}
		}
	}
	return nil
}
