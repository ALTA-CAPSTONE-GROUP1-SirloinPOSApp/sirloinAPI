package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

var (
	JWT_KEY           string = ""
	KEYID             string = ""
	ACCESSKEY         string = ""
	MIDTRANSSERVERKEY string = ""
	AWS_REGION        string = ""
	S3_KEY            string = ""
	S3_SECRET         string = ""
	AWS_BUCKET        string = ""
)

type AppConfig struct {
	DBUser string
	DBPass string
	DBHost string
	DBPort int
	DBName string
	jwtKey string
	keyid  string
	// accesskey         string
	midtransserverkey string
	AWSREGION         string
	S3KEY             string
	S3SECRET          string
	AWSBUCKET         string
}

func InitConfig() *AppConfig {
	return ReadEnv()
}

func ReadEnv() *AppConfig {
	app := AppConfig{}
	isRead := true
	// AWS S3 Bucket
	if val, found := os.LookupEnv("AWS_REGION"); found {
		app.AWSREGION = val
		isRead = false
	}
	if val, found := os.LookupEnv("S3_KEY"); found {
		app.S3KEY = val
		isRead = false
	}
	if val, found := os.LookupEnv("S3_SECRET"); found {
		app.S3SECRET = val
		isRead = false
	}
	if val, found := os.LookupEnv("AWS_BUCKET"); found {
		app.AWSBUCKET = val
		isRead = false
	}

	// midtrans
	if val, found := os.LookupEnv("MIDTRANSSERVERKEY"); found {
		app.keyid = val
		isRead = false
		MIDTRANSSERVERKEY = val
	}

	// JWT
	if val, found := os.LookupEnv("JWT_KEY"); found {
		app.jwtKey = val
		isRead = false
		JWT_KEY = val
	}

	// DATABASE
	if val, found := os.LookupEnv("DBUSER"); found {
		app.DBUser = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBPASS"); found {
		app.DBPass = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBHOST"); found {
		app.DBHost = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBPORT"); found {
		cnv, _ := strconv.Atoi(val)
		app.DBPort = cnv
		isRead = false
	}
	if val, found := os.LookupEnv("DBNAME"); found {
		app.DBName = val
		isRead = false
	}

	if isRead {
		err := godotenv.Load("local.env")
		if err != nil {
			fmt.Println("Error saat baca env", err.Error())
			return nil
		}

		app.DBUser = os.Getenv("DBUSER")
		app.DBPass = os.Getenv("DBPASS")
		app.DBHost = os.Getenv("DBHOST")
		readData := os.Getenv("DBPORT")
		app.DBPort, err = strconv.Atoi(readData)
		if err != nil {
			fmt.Println("Error saat convert", err.Error())
			return nil
		}
		app.DBName = os.Getenv("DBNAME")
		app.jwtKey = os.Getenv("JWTKEY")
		app.midtransserverkey = os.Getenv("MIDTRANSSERVERKEY")
		app.AWSREGION = os.Getenv("AWSREGION")
		app.S3KEY = os.Getenv("S3KEY")
		app.S3SECRET = os.Getenv("S3SECRET")
		app.AWSBUCKET = os.Getenv("AWSBUCKET")

	}

	JWT_KEY = app.jwtKey
	MIDTRANSSERVERKEY = app.midtransserverkey
	AWS_REGION = app.AWSREGION
	S3_KEY = app.S3KEY
	S3_SECRET = app.S3SECRET
	AWS_BUCKET = app.AWSBUCKET

	return &app
}

func MidtransSnapClient() snap.Client {
	s := snap.Client{}
	s.New(MIDTRANSSERVERKEY, midtrans.Sandbox)
	return s
}

func MidtransCoreAPIClient() coreapi.Client {
	c := coreapi.Client{}
	c.New(MIDTRANSSERVERKEY, midtrans.Sandbox)
	return c
}
