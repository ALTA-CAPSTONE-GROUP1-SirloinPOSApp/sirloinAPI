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

	if strings.Contains(msg, "already exist") {
		words := strings.Split(msg, ": ")
		log.Println("error running " + words[0] + " service: already exist")
		resp["message"] = words[0] + " already exist"
		code = http.StatusConflict

	} else if strings.Contains(msg, "secure_password") {
		log.Println("error running register service: the password does not meet security requirements")
		code = http.StatusBadRequest
		resp["message"] = "password must be at least 8 characters long, must contain uppercase letters, must contain lowercase letters, must contain numbers, must not be too general"

	} else if strings.Contains(msg, "required") {
		if strings.Contains(msg, "'Upc'") {
			log.Println("error running register service: required fields")
			code = http.StatusBadRequest
			resp["message"] = "upc shouldn't be empty"
		} else if strings.Contains(msg, "'Category'") {
			log.Println("error running register service: required fields")
			code = http.StatusBadRequest
			resp["message"] = "category shouldn't be empty"
		} else if strings.Contains(msg, "'ProductName'") {
			log.Println("error running register service: required fields")
			code = http.StatusBadRequest
			resp["message"] = "product name shouldn't be empty"
		} else if strings.Contains(msg, "'Stock'") {
			log.Println("error running register service: required fields")
			code = http.StatusBadRequest
			resp["message"] = "stock shouldn't be empty"
		} else if strings.Contains(msg, "'Price'") {
			log.Println("error running register service: required fields")
			code = http.StatusBadRequest
			resp["message"] = "price shouldn't be empty"
		} else {
			log.Println("error running register service: required fields")
			code = http.StatusBadRequest
			resp["message"] = "required fields must be filled"
		}

	} else if strings.Contains(msg, "upc_length") {
		log.Println("error running register service: upc length not meet requirement")
		code = http.StatusBadRequest
		resp["message"] = "upc should only 12 or 13 digits"

	} else if strings.Contains(msg, "'numeric'") {
		logMsg := ""
		if strings.Contains(msg, "RegisterValidate.PhoneNumber") {
			logMsg = "register business phone number"
		} else if strings.Contains(msg, "AddCustomerValidate.PhoneNumber") {
			logMsg = "register customer phone number"
		} else if strings.Contains(msg, "ProductValidate.Upc") {
			logMsg = "register product upc"
		} else if strings.Contains(msg, "ProductValidate.Stock") {
			logMsg = "register product stock"
		} else if strings.Contains(msg, "ProductValidate.MinimumStock") {
			logMsg = "register product minimum stock"
		} else if strings.Contains(msg, "ProductValidate.BuyingPrice") {
			logMsg = "register product buying price"
		} else if strings.Contains(msg, "ProductValidate.Price") {
			logMsg = "register product price"
		} else {
			words := strings.Split(msg, ": ")
			logMsg = words[0]
		}
		log.Println("error running " + logMsg + " service: must be numeric")
		code = http.StatusBadRequest
		resp["message"] = logMsg + " must be a number"

	} else if strings.Contains(msg, "'alpha_space'") {
		logMsg := ""
		if strings.Contains(msg, "RegisterValidate.BusinessName") {
			logMsg = "register business name"
		} else if strings.Contains(msg, "AddCustomerValidate.Name") {
			logMsg = "register customer name"
		} else if strings.Contains(msg, "ProductValidate.Category") {
			logMsg = "register product category"
		} else {
			words := strings.Split(msg, ": ")
			logMsg = words[0]
		}
		log.Println("error running " + logMsg + " service: must be alpha_space")
		code = http.StatusBadRequest
		resp["message"] = logMsg + " are only allowed to contain letters and spaces"

	} else if strings.Contains(msg, "'alpha_space_numeric'") {
		logMsg := ""
		if strings.Contains(msg, "ProductValidate.ProductName") {
			logMsg = "register product name"
		} else {
			words := strings.Split(msg, ": ")
			logMsg = words[0]
		}
		log.Println("error running " + logMsg + " service: must be alpha_space_numeric")
		code = http.StatusBadRequest
		resp["message"] = logMsg + " are only allowed to contain letters, spaces and numeric"

	} else if strings.Contains(msg, "email") {
		logMsg := ""
		if strings.Contains(msg, "RegisterValidate.Email") {
			logMsg = "register business email"
		} else if strings.Contains(msg, "AddCustomerValidate.Email") {
			logMsg = "register customer email"
		} else {
			words := strings.Split(msg, ": ")
			logMsg = words[0]
		}
		log.Println("error running " + logMsg + " service: must be email format")
		code = http.StatusBadRequest
		resp["message"] = "incorrect " + logMsg + " format"

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
