package main

import (
	"log"
	"sirloinapi/config"
	"sirloinapi/migration"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	pd "sirloinapi/features/product/data"
	ph "sirloinapi/features/product/handler"
	ps "sirloinapi/features/product/services"

	ud "sirloinapi/features/user/data"
	uh "sirloinapi/features/user/handler"
	us "sirloinapi/features/user/services"

	td "sirloinapi/features/transaction/data"
	th "sirloinapi/features/transaction/handler"
	ts "sirloinapi/features/transaction/services"

	cd "sirloinapi/features/customer/data"
	ch "sirloinapi/features/customer/handler"
	cs "sirloinapi/features/customer/services"
)

func main() {
	e := echo.New()
	cfg := config.InitConfig()
	db := config.InitDB(*cfg)
	migration.Migrate(db)

	userData := ud.New(db)
	userSrv := us.New(userData)
	userHdl := uh.New(userSrv)

	prodData := pd.New(db)
	prodSrv := ps.New(prodData)
	prodHdl := ph.New(prodSrv)

	transData := td.New(db)
	transSrv := ts.New(transData)
	transHdl := th.New(transSrv)

	cusData := cd.New(db)
	cusSrv := cs.New(cusData)
	cusHdl := ch.New(cusSrv)

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "- method=${method}, uri=${uri}, status=${status}, error=${error}\n",
	}))

	//user
	e.POST("/register", userHdl.Register())
	e.POST("/login", userHdl.Login())
	e.GET("/users", userHdl.Profile(), echojwt.JWT([]byte(config.JWT_KEY)))
	e.PUT("/users", userHdl.Update(), echojwt.JWT([]byte(config.JWT_KEY)))
	e.DELETE("/users", userHdl.Delete(), echojwt.JWT([]byte(config.JWT_KEY)))

	//product
	e.POST("/products", prodHdl.Add(), echojwt.JWT([]byte(config.JWT_KEY)))
	e.GET("/products", prodHdl.GetUserProducts(), echojwt.JWT([]byte(config.JWT_KEY)))
	e.PUT("/products/:product_id", prodHdl.Update(), echojwt.JWT([]byte(config.JWT_KEY)))
	e.DELETE("/products/:product_id", prodHdl.Delete(), echojwt.JWT([]byte(config.JWT_KEY)))
	e.GET("/products/:product_id", prodHdl.GetProductById(), echojwt.JWT([]byte(config.JWT_KEY)))
	e.GET("/products/admin", prodHdl.GetAdminProducts(), echojwt.JWT([]byte(config.JWT_KEY)))

	//customer
	e.POST("/customers", cusHdl.Add(), echojwt.JWT([]byte(config.JWT_KEY)))
	e.GET("/customers", cusHdl.GetUserCustomers(), echojwt.JWT([]byte(config.JWT_KEY)))
	e.GET("/customers/:customer_id", cusHdl.GetCustomerById(), echojwt.JWT([]byte(config.JWT_KEY)))
	e.PUT("/customers/:customer_id", cusHdl.Update(), echojwt.JWT([]byte(config.JWT_KEY)))

	//transaction
	e.POST("/transactions", transHdl.AddSell(), echojwt.JWT([]byte(config.JWT_KEY)))
	e.POST("/transactions/buy", transHdl.AddBuy(), echojwt.JWT([]byte(config.JWT_KEY)))
	e.GET("/transactions", transHdl.GetTransactionHistory(), echojwt.JWT([]byte(config.JWT_KEY)))
	e.GET("/transactions/:transaction_id", transHdl.GetTransactionDetails(), echojwt.JWT([]byte(config.JWT_KEY)))

	if err := e.Start(":8000"); err != nil {
		log.Println(err.Error())
	}
}
