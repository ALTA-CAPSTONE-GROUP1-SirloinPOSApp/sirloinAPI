package helper

import (
	"log"
	"regexp"
	"sirloinapi/features/user"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// user struct validation
type RegisterValidate struct {
	BusinessName string `validate:"required,alpha_space"`
	Email        string `validate:"required,email"`
	Address      string `validate:"required"`
	PhoneNumber  string `validate:"required,numeric"`
	Password     string `validate:"required,secure_password"`
}

type PasswordValidate struct {
	Password string `validate:"secure_password"`
}

func ToValidate(option string, data user.Core) interface{} {
	switch option {
	case "register":
		res := RegisterValidate{}
		res.Email = data.Email
		res.BusinessName = data.BusinessName
		res.PhoneNumber = data.PhoneNumber
		res.Password = data.Password
		res.Address = data.Address
		return res
	case "password":
		res := PasswordValidate{}
		res.Password = data.Password
		return res
	default:
		return nil
	}
}

func alphaSpace(fl validator.FieldLevel) bool {
	match, _ := regexp.MatchString("^[a-zA-Z\\s]+$", fl.Field().String())
	return match
}

func securePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return false
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return false
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return false
	}
	if regexp.MustCompile(`^(?i)(password|1234|qwerty)`).MatchString(password) {
		return false
	}
	return true
}

func Validasi(data interface{}) error {
	validate = validator.New()
	validate.RegisterValidation("alpha_space", alphaSpace)
	validate.RegisterValidation("secure_password", securePassword)
	err := validate.Struct(data)
	if err != nil {
		log.Println("log on helper validasi: ", err)
		return err
	}
	return nil
}
