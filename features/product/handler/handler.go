package handler

import (
	"log"
	"mime/multipart"
	"net/http"
	"sirloinapi/features/product"
	"sirloinapi/helper"
	"strconv"

	"github.com/labstack/echo/v4"
)

type productControl struct {
	srv product.ProductService
}

func New(srv product.ProductService) product.ProductHandler {
	return &productControl{
		srv: srv,
	}
}

func (pc *productControl) Add() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")
		input := AddProductReq{}
		var productImage *multipart.FileHeader

		if err := c.Bind(&input); err != nil {
			log.Println("\tbind input error: ", err.Error())
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("wrong input"))
		}

		file, err := c.FormFile("product_image")
		if file != nil && err == nil {
			productImage = file
		} else if file != nil && err != nil {
			log.Println("\terror read product image: ", err.Error())
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("wrong image input"))
		}

		res, err := pc.srv.Add(token, *ToCore(input), productImage)
		if err != nil {
			log.Println("\terror running add product service: ", err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("server problem"))
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"data":    res,
			"message": "success add product",
		})
	}
}
func (pc *productControl) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")
		var prodImg *multipart.FileHeader

		productId := c.Param("product_id")
		cProdId, _ := strconv.Atoi(productId)

		input := AddProductReq{}
		err := c.Bind(&input)
		if err != nil {
			log.Println("bind input error: ", err.Error())
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("wrong input"))
		}

		file, err := c.FormFile("product_image")
		if file != nil && err == nil {
			prodImg = file
		} else if file != nil && err != nil {
			log.Println("\terror read product image: ", err.Error())
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("wrong image input"))
		}

		res, err := pc.srv.Update(token, uint(cProdId), *ToCore(input), prodImg)
		if err != nil {
			log.Println("\terror running update post service")
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("server problem"))
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"data":    res,
			"message": "success update product",
		})
	}
}
func (pc *productControl) Delete() echo.HandlerFunc {
	return nil
}

func (pc *productControl) GetUserProducts() echo.HandlerFunc {
	return nil
}
func (pc *productControl) GetProductById() echo.HandlerFunc {
	return nil
}
