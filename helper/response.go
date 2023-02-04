package helper

import (
	"log"
	"net/http"
	"strings"
)

func ErrorResponse(msg string) interface{} {
	resp := map[string]interface{}{}
	resp["message"] = msg

	return resp
}

func PrintErrorResponse(msg string) (int, interface{}) {
	resp := map[string]interface{}{}
	code := -1
	if msg != "" {
		resp["message"] = msg
	}

	if strings.Contains(msg, "user already exist") {
		log.Println("error running register service: user already exist")
		resp["message"] = "email already exist"
		code = http.StatusConflict
	} else if strings.Contains(msg, "phone number already exist") {
		log.Println("error running register service: phone number already exist")
		code = http.StatusConflict
		resp["message"] = "phone number already exist"
	} else if strings.Contains(msg, "secure_password") {
		log.Println("error running register service: the password does not meet security requirements")
		code = http.StatusBadRequest
		resp["message"] = "password must be at least 8 characters long, must contain uppercase letters, must contain lowercase letters, must contain numbers, must not be too general"
	} else if strings.Contains(msg, "required") {
		log.Println("error running register service: required fields")
		code = http.StatusBadRequest
		resp["message"] = "required fields must be filled"
	} else if strings.Contains(msg, "PhoneNumber") && strings.Contains(msg, "numeric") {
		log.Println("error running register service: phone number must be numeric")
		code = http.StatusBadRequest
		resp["message"] = "the phone number must be a number"
	} else if strings.Contains(msg, "BusinessName") && strings.Contains(msg, "alpha_space") {
		log.Println("error running register service: business names must be alpha_space")
		code = http.StatusBadRequest
		resp["message"] = "business names are only allowed to contain letters and spaces"
	} else if strings.Contains(msg, "Email") && strings.Contains(msg, "email") {
		log.Println("error running register service: Email must be email format")
		code = http.StatusBadRequest
		resp["message"] = "incorrect e-mail format"
	} else if strings.Contains(msg, "token error") && strings.Contains(msg, "customer") {
		log.Println("error running register service: extract token error")
		code = http.StatusBadRequest
		resp["message"] = "extract token error/not allowed to access customer"
	} else {
		log.Println("error running register service")
		code = http.StatusInternalServerError
		resp["message"] = "server problem"
	}

	return code, resp
}
