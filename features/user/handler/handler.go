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
			return c.JSON(helper.PrintErrorResponse(err.Error()))
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

func (uc *userControl) Profile() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		res, err := uc.srv.Profile(token)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Println("user not found: ", err.Error())
				return c.JSON(http.StatusNotFound, helper.ErrorResponse("user not found"))
			} else {
				log.Println("error profile service: ", err.Error())
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("server problem"))
			}
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"data":    ToResponse(res),
			"message": "success get tenant profile",
		})
	}
}

func (uc *userControl) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		updatedData := RegisterRequest{}
		if err := c.Bind(&updatedData); err != nil {
			return c.JSON(http.StatusBadRequest, "wrong input format")
		}

		res, err := uc.srv.Update(token, *ToCore(updatedData))
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Println("user not found: ", err.Error())
				return c.JSON(http.StatusNotFound, helper.ErrorResponse("user not found"))
			} else {
				log.Println("error update service: ", err.Error())
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("server problem"))
			}
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"data":    res,
			"message": "success update user's data",
		})
	}
}

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
