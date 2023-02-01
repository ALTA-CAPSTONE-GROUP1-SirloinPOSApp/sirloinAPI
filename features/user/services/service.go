package service

import (
	"errors"
	"sirloinapi/features/user"
	"sirloinapi/helper"
	"strings"

	"github.com/go-playground/validator/v10"
)

type userUseCase struct {
	qry user.UserData
	vld *validator.Validate
}

func New(ud user.UserData) user.UserService {
	return &userUseCase{
		qry: ud,
		vld: validator.New(),
	}
}

func (uuc *userUseCase) Register(newUser user.Core) (user.Core, error) {
	err := helper.Validasi(helper.ToRegister(newUser))
	if err != nil {
		return user.Core{}, err
	}
	hashed := helper.GeneratePassword(newUser.Password)
	newUser.Password = hashed

	res, err := uuc.qry.Register(newUser)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "Duplicate") && strings.Contains(err.Error(), "users.email") {
			msg = "user already exist"
		} else if strings.Contains(err.Error(), "Duplicate") && strings.Contains(err.Error(), "users.phone_number") {
			msg = "phone number already exist"
		} else {
			msg = "server problem"
		}
		return user.Core{}, errors.New(msg)
	}

	return res, nil
}

// func (uuc *userUseCase) Login(email, password string) (string, user.Core, error) {
// 	res, err := uuc.qry.Login(email)

// 	if err != nil {
// 		errmsg := ""
// 		if strings.Contains(err.Error(), "not found") {
// 			errmsg = err.Error()
// 		} else {
// 			errmsg = "server problem"
// 		}
// 		log.Println("error login query: ", err.Error())
// 		return "", user.Core{}, errors.New(errmsg)
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(password)); err != nil {
// 		log.Println("wrong password :", err.Error())
// 		return "", user.Core{}, errors.New("wrong password")
// 	}

// 	claims := jwt.MapClaims{}
// 	claims["authorized"] = true
// 	claims["userID"] = res.ID
// 	// claims["exp"] = time.Now().Add(time.Hour * 1).Unix() //Token expires after 1 hour
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	useToken, _ := token.SignedString([]byte(config.JWT_KEY))

// 	return useToken, res, nil
// }

// func (uuc *userUseCase) Profile(userToken interface{}) (user.Core, error) {
// 	id := helper.ExtractToken(userToken)
// 	if id <= 0 {
// 		log.Println("error extraxt token")
// 		return user.Core{}, errors.New("data not found")
// 	}
// 	res, err := uuc.qry.Profile(uint(id))
// 	if err != nil {
// 		errmsg := ""
// 		if strings.Contains(err.Error(), "not found") {
// 			errmsg = "data not found"
// 		} else {
// 			errmsg = "server problem"
// 		}
// 		log.Println("error profile query: ", err.Error())
// 		return user.Core{}, errors.New(errmsg)
// 	}
// 	return res, nil
// }

// func (uuc *userUseCase) Update(userToken interface{}, updateData user.Core) (user.Core, error) {
// 	userId := helper.ExtractToken(userToken)
// 	if userId <= 0 {
// 		log.Println("extract token error")
// 		return user.Core{}, errors.New("extract token error")
// 	}
// 	if updateData.Password != "" {
// 		hashed, _ := bcrypt.GenerateFromPassword([]byte(updateData.Password), bcrypt.DefaultCost)
// 		updateData.Password = string(hashed)
// 	}

// 	res, err := uuc.qry.Update(uint(userId), updateData)
// 	if err != nil {
// 		errmsg := ""
// 		if strings.Contains(err.Error(), "not found") {
// 			errmsg = "data not found"
// 		} else {
// 			errmsg = "server problem"
// 		}
// 		log.Println("error update query: ", err.Error())
// 		return user.Core{}, errors.New(errmsg)
// 	}
// 	return res, nil
// }

// func (uuc *userUseCase) Delete(userToken interface{}) error {
// 	userId := helper.ExtractToken(userToken)
// 	if userId <= 0 {
// 		return errors.New("data not found")
// 	}
// 	err := uuc.qry.Delete(uint(userId))
// 	if err != nil {
// 		msg := ""
// 		if strings.Contains(err.Error(), "not found") {
// 			msg = "data not found"
// 		} else {
// 			msg = "server problem"
// 		}
// 		return errors.New(msg)
// 	}
// 	return nil
// }
