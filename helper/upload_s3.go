package helper

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	_config "sirloinapi/config"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadImageToS3(fileName string, fileData multipart.File) (string, error) {
	// The session the S3 Uploader will use
	sess := _config.GetSession()

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(_config.AWS_BUCKET),
		Key:         aws.String(fileName),
		Body:        fileData,
		ContentType: aws.String("image"),
		ACL:         aws.String("public-read"),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file, %v", err)
	}
	return result.Location, nil
}

func UploadPdfToS3(fileName string, fileData multipart.File) (string, error) {
	// The session the S3 Uploader will use
	sess := _config.GetSession()

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(_config.AWS_BUCKET),
		Key:         aws.String(fileName),
		Body:        fileData,
		ContentType: aws.String("application/pdf"),
		ACL:         aws.String("public-read"),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file, %v", err)
	}
	return result.Location, nil
}

func CheckFileExtensionImage(filename string) (string, error) {
	extension := strings.ToLower(filename[strings.LastIndex(filename, ".")+1:])

	if extension != "jpg" && extension != "jpeg" && extension != "png" {
		return "", fmt.Errorf("forbidden file type")
	}
	return extension, nil
}

func CheckFileExtensionPdf(filename string) (string, error) {
	extension := strings.ToLower(filename[strings.LastIndex(filename, ".")+1:])

	if extension != "pdf" {
		return "", fmt.Errorf("forbidden file type")
	}
	return extension, nil
}

func CheckFileSize(size int64) error {
	var fileSize int64 = 1097152
	if size == 0 {
		return fmt.Errorf("illegal file size")
	}

	if size > fileSize {
		return fmt.Errorf("file size too big, %d MB", fileSize/1000000)
	}

	return nil
}

func UploadProductPhotoS3(file multipart.FileHeader, productId int) (string, error) {
	s3Session := _config.GetSession()
	uploader := s3manager.NewUploader(s3Session)
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	ext := filepath.Ext(file.Filename)

	cnv := strconv.Itoa(productId)
	res, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(_config.AWS_BUCKET),
		Key:    aws.String("files/products/" + cnv + "/product_photo_" + fmt.Sprint(productId) + ext),
		Body:   src,
	})
	if err != nil {
		return "", errors.New("problem with upload post photo")
	}

	return res.Location, nil
}
