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

	existed := 0
	pq.db.Raw("SELECT COUNT(*) FROM products p WHERE p.product_name = ?", cnvP.ProductName).Scan(&existed)
	if existed >= 1 {
		log.Println("duplicated product")
		return product.Core{}, errors.New("duplicated product")
	}
	existed = 0
	pq.db.Raw("SELECT COUNT(*) FROM products p WHERE p.upc = ?", cnvP.Upc).Scan(&existed)
	if existed >= 1 {
		log.Println("duplicated product")
		return product.Core{}, errors.New("duplicated product")
	}

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
	cnvP := CoreToData(updProduct)
	cnvP.UserId = userId

	if productImage != nil {
		path, err := helper.UploadProductPhotoS3(*productImage, int(productId))
		if err != nil {
			log.Println("\terror upload product photo: ", err.Error())
			return product.Core{}, err
		}
		cnvP.ProductImage = path
	}

	qry := pq.db.Where("id = ? AND user_id = ?", productId, userId).Updates(&cnvP)
	if qry.RowsAffected <= 0 {
		log.Println("\tupdate product query error: data not found")
		return product.Core{}, errors.New("not found")
	}

	if err := qry.Error; err != nil {
		log.Println("\tupdate product query error: ", err.Error())
		return product.Core{}, errors.New("not found")
	}

	return DataToCore(cnvP), nil
}
func (pq *productQuery) Delete(userId, productId uint) error {
	qry := pq.db.Where("user_id = ?", userId).Delete(&Product{}, productId)

	if aff := qry.RowsAffected; aff <= 0 {
		log.Println("\tno rows affected: data not found")
		return errors.New("data not found")
	}

	if err := qry.Error; err != nil {
		log.Println("\tdelete query error: ", err.Error())
		return err
	}

	return nil
}
func (pq *productQuery) GetUserProducts(userId uint) ([]product.Core, error) {
	userProd := []product.Core{}
	err := pq.db.Raw("SELECT p.id , upc , category , product_name , minimum_stock , stock , buying_price , price , product_image , supplier , items_sold FROM products p JOIN users u ON u.id = p.user_id WHERE p.deleted_at IS NULL AND user_id = ?", userId).Scan(&userProd).Error
	if err != nil {
		log.Println("\terror query get user product: ", err.Error())
		return []product.Core{}, err
	}

	return userProd, nil
}
func (pq *productQuery) GetProductById(userId, productId uint) (product.Core, error) {
	prod := product.Core{}
	err := pq.db.Raw("SELECT p.id , upc , category , product_name , minimum_stock , stock , buying_price , price , product_image , supplier , items_sold FROM products p JOIN users u ON u.id = p.user_id WHERE p.deleted_at IS NULL AND p.id = ? AND u.id = ?", productId, userId).Scan(&prod).Error
	if err != nil {
		log.Println("\terror query get all product: ", err.Error())
		return product.Core{}, err
	}

	return prod, nil
}
func (pq *productQuery) GetAdminProducts() ([]product.Core, error) {
	userProd := []product.Core{}
	err := pq.db.Raw("SELECT p.id , upc , category , product_name , minimum_stock , stock , buying_price , price , product_image , supplier , items_sold FROM products p JOIN users u ON u.id = p.user_id WHERE p.deleted_at IS NULL AND user_id = 1").Scan(&userProd).Error
	if err != nil {
		log.Println("\terror query get user product: ", err.Error())
		return []product.Core{}, err
	}

	return userProd, nil
}
