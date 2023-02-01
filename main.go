package main

import (
	"log"
	"sirloinapi/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	ud "sirloinapi/features/user/data"
	uh "sirloinapi/features/user/handler"
	us "sirloinapi/features/user/services"
)

func main() {
	e := echo.New()
	cfg := config.InitConfig()
	db := config.InitDB(*cfg)
	config.Migrate(db)

	userData := ud.New(db)
	userSrv := us.New(userData)
	userHdl := uh.New(userSrv)

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "- method=${method}, uri=${uri}, status=${status}, error=${error}\n",
	}))

	//user
	e.POST("/register", userHdl.Register())

	if err := e.Start(":8000"); err != nil {
		log.Println(err.Error())
	}
}
