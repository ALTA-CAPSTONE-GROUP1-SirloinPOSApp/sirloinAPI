package user

import (
	"github.com/labstack/echo/v4"
)

type Core struct {
	ID           uint
	BusinessName string
	Email        string
	Address      string
	PhoneNumber  string
	Password     string
}

type UserHandler interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	// Profile() echo.HandlerFunc
	// Update() echo.HandlerFunc
	// Delete() echo.HandlerFunc
}

type UserService interface {
	Register(newUser Core) (Core, error)
	Login(email, password string) (string, Core, error)
	// Profile(userToken interface{}) (Core, error)
	// Update(userToken interface{}, updateData Core) (Core, error)
	// Delete(userToken interface{}) error
}

type UserData interface {
	Register(newUser Core) (Core, error)
	Login(email string) (Core, error)
	Profile(id uint) (Core, error)
	// Update(id uint, updateData Core) (Core, error)
	// Delete(id uint) error
}
