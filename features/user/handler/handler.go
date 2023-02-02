package handler

import (
	"log"
	"net/http"
	"sirloinapi/features/user"
	"sirloinapi/helper"
	"strings"

	"github.com/labstack/echo/v4"
)

type userControl struct {
	srv user.UserService
}

func New(srv user.UserService) user.UserHandler {
	return &userControl{
		srv: srv,
	}
}

func (uc *userControl) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := RegisterRequest{}

		if err := c.Bind(&input); err != nil {
			log.Println("error bind input")
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("wrong input"))
		}

		res, err := uc.srv.Register(*ToCore(input))
		if err != nil {
			log.Println("log on user handler: ", err)
			if strings.Contains(err.Error(), "user already exist") {
				log.Println("error running register service: user already exist")
				return c.JSON(http.StatusConflict, helper.ErrorResponse("user or email already exist"))
			} else if strings.Contains(err.Error(), "phone number already exist") {
				log.Println("error running register service: phone number already exist")
				return c.JSON(http.StatusConflict, helper.ErrorResponse("phone number already exist"))
			} else if strings.Contains(err.Error(), "secure_password") {
				log.Println("error running register service: the password does not meet security requirements")
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse("password must be at least 8 characters long, must contain uppercase letters, must contain lowercase letters, must contain numbers, must not be too general"))
			} else if strings.Contains(err.Error(), "required") {
				log.Println("error running register service: required fields")
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse("required fields must be filled"))
			} else if strings.Contains(err.Error(), "PhoneNumber") && strings.Contains(err.Error(), "numeric") {
				log.Println("error running register service: phone number must be numeric")
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse("the phone number must be a number"))
			} else if strings.Contains(err.Error(), "BusinessName") && strings.Contains(err.Error(), "alpha_space") {
				log.Println("error running register service: business names must be alpha_space")
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse("business names are only allowed to contain letters and spaces"))
			} else {
				log.Println("error running register service")
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("server problem"))
			}
		}
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"data":    ToResponse(res),
			"message": "success register",
		})
	}
}

func (uc *userControl) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := LoginReqest{}
		if err := c.Bind(&input); err != nil {
			log.Println("error login request: ", err.Error())
			return c.JSON(http.StatusBadRequest, "wrong input")
		}

		token, res, err := uc.srv.Login(input.Email, input.Password)
		if err != nil {
			if strings.Contains(err.Error(), "password") {
				log.Println("wrong password: ", err.Error())
				return c.JSON(http.StatusUnauthorized, helper.ErrorResponse("wrong password"))
			} else if strings.Contains(err.Error(), "not found") {
				log.Println("user not found: ", err.Error())
				return c.JSON(http.StatusNotFound, helper.ErrorResponse("wrong email"))
			} else {
				log.Println("error login service: ", err.Error())
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("server problem"))
			}
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"data":    ToLoginResp(res, token),
			"message": "login success",
		})
	}
}

// func (uc *userControl) Profile() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		token := c.Get("user")

// 		res, err := uc.srv.Profile(token)
// 		if err != nil {
// 			if strings.Contains(err.Error(), "not found") {
// 				log.Println("user not found: ", err.Error())
// 				return c.JSON(http.StatusNotFound, helper.ErrorResponse("user not found"))
// 			} else {
// 				log.Println("error profile service: ", err.Error())
// 				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("server problem"))
// 			}
// 		}

// 		return c.JSON(http.StatusOK, map[string]interface{}{
// 			"data":    ToResponse(res),
// 			"message": "get profile success",
// 		})
// 	}
// }

// func (uc *userControl) Update() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		token := c.Get("user")

// 		updatedData := RegisterRequest{}
// 		if err := c.Bind(&updatedData); err != nil {
// 			return c.JSON(http.StatusBadRequest, "wrong input format")
// 		}

// 		res, err := uc.srv.Update(token, *ToCore(updatedData))
// 		if err != nil {
// 			if strings.Contains(err.Error(), "not found") {
// 				log.Println("user not found: ", err.Error())
// 				return c.JSON(http.StatusNotFound, helper.ErrorResponse("user not found"))
// 			} else {
// 				log.Println("error update service: ", err.Error())
// 				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("server problem"))
// 			}
// 		}

// 		return c.JSON(http.StatusOK, map[string]interface{}{
// 			"data":    res,
// 			"message": "success update user's data",
// 		})
// 	}
// }

// func (uc *userControl) Delete() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		token := c.Get("user")
// 		err := uc.srv.Delete(token)
// 		if err != nil {
// 			if strings.Contains(err.Error(), "not found") {
// 				c.JSON(http.StatusNotFound, helper.ErrorResponse("user not found"))
// 			} else {
// 				c.JSON(http.StatusInternalServerError, helper.ErrorResponse("server problem"))
// 			}
// 		}
// 		return c.JSON(http.StatusOK, map[string]interface{}{
// 			"message": "success delete user",
// 		})
// 	}
// }
