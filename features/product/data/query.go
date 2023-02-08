package data

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"sirloinapi/features/product"
	"sirloinapi/helper"
	"strings"

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
	pq.db.Raw("SELECT COUNT(*) FROM products p WHERE p.product_name = ? AND user_id = ?", cnvP.ProductName, userId).Scan(&existed)
	if existed >= 1 {
		log.Println("duplicated product")
		return product.Core{}, errors.New("duplicated product")
	}
	existed = 0
	pq.db.Raw("SELECT COUNT(*) FROM products p WHERE p.upc = ? AND user_id = ?", cnvP.Upc, userId).Scan(&existed)
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
		src, err := productImage.Open()
		if err != nil {
			return product.Core{}, errors.New("format input file tidak dapat dibuka")
		}
		err = helper.CheckFileSize(productImage.Size)
		if err != nil {
			idx := strings.Index(err.Error(), ",")
			msg := err.Error()
			return product.Core{}, errors.New("format input file size tidak diizinkan, size melebihi" + msg[idx+1:])
		}
		extension, err := helper.CheckFileExtensionImage(productImage.Filename)
		if err != nil {
			return product.Core{}, errors.New("format input file type tidak diizinkan")
		}
		filename := "files/products/" + fmt.Sprint(cnvP.ID) + "/product_photo_" + fmt.Sprint(cnvP.ID) + "." + extension

		path, err := helper.UploadImageToS3(filename, src)
		if err != nil {
			log.Println(errors.New("upload to s3 bucket failed"))
			return product.Core{}, errors.New("upload to s3 bucket failed")
		}
		if len(path) > 0 {
			qry := pq.db.First(&cnvP)
			cnvP.ProductImage = path
			qry.Save(&cnvP)
		}
		defer src.Close()
	}

	return DataToCore(cnvP), nil
}
func (pq *productQuery) Update(userId, productId uint, updProduct product.Core, productImage *multipart.FileHeader) (product.Core, error) {
	cnvP := CoreToData(updProduct)
	cnvP.UserId = userId

	if productImage != nil {
		src, err := productImage.Open()
		if err != nil {
			return product.Core{}, errors.New("format input file tidak dapat dibuka")
		}
		err = helper.CheckFileSize(productImage.Size)
		if err != nil {
			idx := strings.Index(err.Error(), ",")
			msg := err.Error()
			return product.Core{}, errors.New("format input file size tidak diizinkan, size melebihi" + msg[idx+1:])
		}
		extension, err := helper.CheckFileExtensionImage(productImage.Filename)
		if err != nil {
			return product.Core{}, errors.New("format input file type tidak diizinkan")
		}
		filename := "files/products/" + fmt.Sprint(productId) + "/product_photo_" + fmt.Sprint(productId) + "." + extension

		path, err := helper.UploadImageToS3(filename, src)
		if err != nil {
			log.Println(errors.New("upload to s3 bucket failed"))
			return product.Core{}, errors.New("upload to s3 bucket failed")
		}
		if len(path) > 0 {
			cnvP.ProductImage = path
		}
		defer src.Close()
	} else {
		cnvP.ProductImage = ""
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
func (pq *productQuery) GetUserProducts(userId uint, search string) ([]product.Core, error) {
	userProd := []product.Core{}
	var err error

	if search == "" {
		err = pq.db.Raw("SELECT p.id , upc , category , product_name , minimum_stock , stock , buying_price , price , product_image , supplier , items_sold FROM products p JOIN users u ON u.id = p.user_id WHERE p.deleted_at IS NULL AND user_id = ?", userId).Scan(&userProd).Error
	} else {
		err = pq.db.Raw("SELECT p.id , upc , category , product_name , minimum_stock , stock , buying_price , price , product_image , supplier , items_sold FROM products p JOIN users u ON u.id = p.user_id WHERE p.deleted_at IS NULL AND user_id = ? AND product_name LIKE ?", userId, "%"+search+"%").Scan(&userProd).Error
	}
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
func (pq *productQuery) GetAdminProducts(search string) ([]product.Core, error) {
	userProd := []product.Core{}
	var err error
	if search == "" {
		err = pq.db.Raw("SELECT p.id , upc , category , product_name , minimum_stock , stock , buying_price , price , product_image , supplier , items_sold FROM products p JOIN users u ON u.id = p.user_id WHERE p.deleted_at IS NULL AND user_id = 1").Scan(&userProd).Error
	} else {
		err = pq.db.Raw("SELECT p.id , upc , category , product_name , minimum_stock , stock , buying_price , price , product_image , supplier , items_sold FROM products p JOIN users u ON u.id = p.user_id WHERE p.deleted_at IS NULL AND user_id = 1 AND product_name LIKE ?", "%"+search+"%").Scan(&userProd).Error
	}
	if err != nil {
		log.Println("\terror query get user product: ", err.Error())
		return []product.Core{}, err
	}

	return userProd, nil
}
