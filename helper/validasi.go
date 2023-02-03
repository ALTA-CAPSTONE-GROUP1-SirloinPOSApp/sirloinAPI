package helper

import (
	"log"
	"regexp"
	"sirloinapi/features/customer"
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

type AddCustomerValidate struct {
	Name        string `validate:"required,alpha_space"`
	Email       string `validate:"required,email"`
	Address     string `validate:"required"`
	PhoneNumber string `validate:"required,numeric"`
}

type PasswordValidate struct {
	Password string `validate:"secure_password"`
}

type EmailValidate struct {
	Email string `validate:"email"`
}

type PhoneNumberValidate struct {
	PhoneNumber string `validate:"numeric"`
}
type BusinessNameValidate struct {
	BusinessName string `validate:"alpha_space"`
}

func ToValidate(option string, data interface{}) interface{} {
	switch option {
	case "register":
		res := RegisterValidate{}
		if v, ok := data.(user.Core); ok {
			res.Email = v.Email
			res.BusinessName = v.BusinessName
			res.PhoneNumber = v.PhoneNumber
			res.Password = v.Password
			res.Address = v.Address
		}
		return res
	case "password":
		res := PasswordValidate{}
		if v, ok := data.(user.Core); ok {
			res.Password = v.Password
		}
		return res
	case "email":
		res := EmailValidate{}
		if v, ok := data.(user.Core); ok {
			res.Email = v.Email
		}
		return res
	case "pn":
		res := PhoneNumberValidate{}
		if v, ok := data.(user.Core); ok {
			res.PhoneNumber = v.PhoneNumber
		}
		return res
	case "bn":
		res := BusinessNameValidate{}
		if v, ok := data.(user.Core); ok {
			res.BusinessName = v.BusinessName
		}
		return res
	case "customer":
		res := AddCustomerValidate{}
		if v, ok := data.(customer.Core); ok {
			res.Email = v.Email
			res.Name = v.Name
			res.PhoneNumber = v.PhoneNumber
			res.Address = v.Address
		}
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
