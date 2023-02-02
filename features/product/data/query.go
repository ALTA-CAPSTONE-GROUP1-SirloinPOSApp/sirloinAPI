package data

import (
	"errors"
	"log"
	"mime/multipart"
	"sirloinapi/features/product"
	"sirloinapi/helper"

	"gorm.io/gorm"
)

type productQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) product.ProductData {
	return &productQuery{
		db: db,
	}
}

func (pq *productQuery) Add(userId uint, newProduct product.Core, productImage *multipart.FileHeader) (product.Core, error) {
	cnvP := CoreToData(newProduct)
	cnvP.UserId = userId

	err := pq.db.Create(&cnvP).Error
	if err != nil {
		log.Println("\tadd product query error: ", err.Error())
		return product.Core{}, errors.New("server problem")
	}

	if productImage != nil {
		path, err := helper.UploadProductPhotoS3(*productImage, int(cnvP.ID))
		if err != nil {
			log.Println("\terror upload product photo: ", err.Error())
			return product.Core{}, err
		}
		qry := pq.db.First(&cnvP)
		cnvP.ProductImage = path
		qry.Save(&cnvP)
	}

	return DataToCore(cnvP), nil
}
func (pq *productQuery) Update(userId, productId uint, updProduct product.Core, productImage *multipart.FileHeader) (product.Core, error) {
	return product.Core{}, nil
}
func (pq *productQuery) Delete(userId, productId uint) error {
	return nil
}
func (pq *productQuery) GetUserProducts(userId uint) ([]product.Core, error) {
	return []product.Core{}, nil
}
func (pq *productQuery) GetProductById(productId uint) (product.Core, error) {

	return product.Core{}, nil
}
