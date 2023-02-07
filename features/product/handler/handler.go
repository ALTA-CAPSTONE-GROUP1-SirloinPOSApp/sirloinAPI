package handler

import (
	"log"
	"mime/multipart"
	"net/http"
	"sirloinapi/features/product"
	"sirloinapi/helper"
	"strconv"
	"strings"

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

		if input.Upc == "" {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("upc shouldn't be empty"))
		} else if input.ProductName == "" {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("product_name shouldn't be empty"))
		} else if input.Category == "" {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("category shouldn't be empty"))
		} else if input.Price == 0.0 {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("price shouldn't be empty"))
		} else if input.Stock == 0 {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("stock shouldn't be empty"))
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
			if strings.Contains(err.Error(), "duplicated") {
				log.Println("\terror running add product service: ", err.Error())
				return c.JSON(http.StatusConflict, helper.ErrorResponse("duplicated product"))
			} else if strings.Contains(err.Error(), "server") {
				log.Println("\terror running add product service: ", err.Error())
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("server problem"))
			} else if strings.Contains(err.Error(), "format") {
				log.Println("\terror running add product service: ", err.Error())
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
			} else {
				log.Println("\terror running add product service: ", err.Error())
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse("bad request: UPC, stock, minimum_stock, buying_price, price should only numeric and product_name, category should only alpha space, "))
			}
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
			if strings.Contains(err.Error(), "duplicated") {
				log.Println("\terror running add product service: ", err.Error())
				return c.JSON(http.StatusConflict, helper.ErrorResponse("duplicated product"))
			} else if strings.Contains(err.Error(), "server") {
				log.Println("\terror running add product service: ", err.Error())
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("server problem"))
			} else if strings.Contains(err.Error(), "format") {
				log.Println("\terror running add product service: ", err.Error())
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
			} else {
				log.Println("\terror running add product service: ", err.Error())
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse("bad request: UPC, stock, minimum_stock, buying_price, price should only numeric and product_name, category should only alpha space, "))
			}
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"data":    res,
			"message": "success update product",
		})
	}
}
func (pc *productControl) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")
		input := c.Param("product_id")
		cnv, err := strconv.Atoi(input)
		if err != nil {
			log.Println("\tRead param error: ", err.Error())
			return c.JSON(http.StatusBadRequest, "wrong product id parameter")
		}

		err = pc.srv.Delete(token, uint(cnv))
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Println("error calling delete product service: ", err.Error())
				return c.JSON(http.StatusNotFound, helper.ErrorResponse("product not found"))
			} else {
				log.Println("error calling delete product service: ", err.Error())
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("server problem"))
			}
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success delete product",
		})
	}
}

func (pc *productControl) GetUserProducts() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")
		search := c.QueryParam("search")

		res, err := pc.srv.GetUserProducts(token, search)
		if err != nil {
			log.Println("error running GetAllProducts service: ", err.Error())
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusNotFound, helper.ErrorResponse("data not found"))
			} else {
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("server problem"))
			}
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"data":    ToGetProdsResp(res),
			"message": "success show all products",
		})
	}
}
func (pc *productControl) GetProductById() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		input := c.Param("product_id")
		cnv, err := strconv.Atoi(input)
		if err != nil {
			log.Println("\tRead param error: ", err.Error())
			return c.JSON(http.StatusBadRequest, "wrong product id parameter")
		}

		res, err := pc.srv.GetProductById(token, uint(cnv))
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Println("error calling delete product service: ", err.Error())
				return c.JSON(http.StatusNotFound, helper.ErrorResponse("product not found"))
			} else {
				log.Println("error calling delete product service: ", err.Error())
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("server problem"))
			}
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"data":    ToGetProdResp(res),
			"message": "success get product by id",
		})
	}
}
func (pc *productControl) GetAdminProducts() echo.HandlerFunc {
	return func(c echo.Context) error {
		search := c.QueryParam("search")
		res, err := pc.srv.GetAdminProducts(search)
		if err != nil {
			log.Println("error running GetAllProducts service: ", err.Error())
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusNotFound, helper.ErrorResponse("data not found"))
			} else {
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("server problem"))
			}
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"data":    ToGetProdsResp(res),
			"message": "success show all products",
		})
	}
}
