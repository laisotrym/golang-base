package safeweb_lib_validate

import (
	"log"
	"regexp"
	"strings"

	"github.com/nyaruka/phonenumbers"
	"gopkg.in/go-playground/validator.v9"
)

// Validate core user
func ValidateUserName(fl validator.FieldLevel) bool {
	name := fl.Field().String()

	if len(name) < 3 {
		return false
	}

	return !strings.Contains(name, " ")
}

func ValidatePhone(fl validator.FieldLevel) bool {
	p, err := phonenumbers.Parse(fl.Field().String(), "US")
	if err != nil {
		log.Println("[ERROR] Phone verification:", err)
		return false
	}
	return phonenumbers.IsPossibleNumber(p)
}

func ValidateStrongPassword(fl validator.FieldLevel) bool {
	ps := fl.Field().String()
	if len(ps) < 8 {
		return false
	}
	num := `[0-9]{1}`
	a_z := `[a-z]{1}`
	A_Z := `[A-Z]{1}`
	symbol := `[!@#~$%^&*()+|_]{1}`
	if b, err := regexp.MatchString(num, ps); !b || err != nil {
		// return fmt.Errorf("password need number :%v", err)
		return false
	}
	if b, err := regexp.MatchString(a_z, ps); !b || err != nil {
		// return fmt.Errorf("password need lowercase character :%v", err)
		return false
	}
	if b, err := regexp.MatchString(A_Z, ps); !b || err != nil {
		// return fmt.Errorf("password need uppercase character :%v", err)
		return false
	}
	if b, err := regexp.MatchString(symbol, ps); !b || err != nil {
		// return fmt.Errorf("password need special character :%v", err)
		return false
	}
	return true
}

func ValidateCaptcha(fl validator.FieldLevel) bool {
	token := fl.Field().String()

	if token == "" {
		return false
	}

	// Validate reCAPTCHA
	ok, err := ValidateReCAPTCHA(token)
	if err != nil {
		log.Println("[ERROR] reCAPTCHA verification:", err)
	}

	return ok
}

func ValidateSimplePassword(fl validator.FieldLevel) bool {
	ps := fl.Field().String()

	if ps == "" {
		return true
	}

	if len(ps) < 4 {
		return false
	}

	return true
}

func ValidateFullName(fl validator.FieldLevel) bool {
	name := fl.Field().String()

	if strings.TrimSpace(name) == "" {
		return false
	}

	return true
}

func ValidateAlias(fl validator.FieldLevel) bool {
	name := fl.Field().String()

	if len(name) < 3 {
		return false
	}

	return !strings.Contains(name, " ")
}
