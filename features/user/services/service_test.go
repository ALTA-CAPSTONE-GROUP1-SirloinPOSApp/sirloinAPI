package service

import (
	"errors"
	"sirloinapi/features/user"
	"sirloinapi/helper"
	"sirloinapi/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	data := mocks.NewUserData(t)
	password := "Amr12345"
	hash := helper.GeneratePassword(password)
	newUser := user.Core{
		Email:        "mfauzanptra@gmail.com",
		BusinessName: "Muhamad Fauzan Putra",
		PhoneNumber:  "085659171799",
		Address:      "Jln. Lembayung No 24, Bantul, Yogyakarta",
	}
	expectedData := user.Core{
		Email:        "mfauzanptra@gmail.com",
		BusinessName: "Muhamad Fauzan Putra",
		PhoneNumber:  "085659171799",
		Address:      "Jln. Lembayung No 24, Bantul, Yogyakarta",
	}

	t.Run("success register", func(t *testing.T) {
		newUser.Password = hash
		data.On("Register", newUser).Return(expectedData, nil).Once()
		srv := New(data)
		newUser.Password = password
		res, err := srv.Register(newUser)
		assert.Nil(t, err)
		assert.Equal(t, expectedData.ID, res.ID)
		assert.Equal(t, expectedData.BusinessName, res.BusinessName)
		data.AssertExpectations(t)
	})

	t.Run("Duplicate email", func(t *testing.T) {
		newUser.Password = hash
		data.On("Register", newUser).Return(user.Core{}, errors.New("Duplicate users.email")).Once()
		srv := New(data)
		newUser.Password = password
		res, err := srv.Register(newUser)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "user already exist")
		assert.Equal(t, res.BusinessName, "")
		data.AssertExpectations(t)

	})

	t.Run("Duplicate phone number", func(t *testing.T) {
		newUser.Password = hash
		data.On("Register", newUser).Return(user.Core{}, errors.New("Duplicate users.phone_number")).Once()
		srv := New(data)
		newUser.Password = password
		res, err := srv.Register(newUser)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "phone number already exist")
		assert.Equal(t, res.BusinessName, "")
		data.AssertExpectations(t)

	})

	t.Run("server problem", func(t *testing.T) {
		newUser.Password = hash
		data.On("Register", newUser).Return(user.Core{}, errors.New("server error")).Once()
		srv := New(data)
		newUser.Password = password
		res, err := srv.Register(newUser)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Equal(t, res.BusinessName, "")
		data.AssertExpectations(t)

	})

	t.Run("validation problem", func(t *testing.T) {
		newUser.Password = hash
		srv := New(data)
		newUser.Password = password
		newUser.PhoneNumber = "ojaosdjoasidj"
		res, err := srv.Register(newUser)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "validation")
		assert.Equal(t, res.BusinessName, "")

	})

}

// func TestLogn(t *testing.T) {
// 	data := mocks.NewUserData(t)
// 	t.Run("succcess login", func(t *testing.T) {
// 		email := "mfauzanptra@gmail.com"
// 		password := "paupau99"
// 		hashed, _ := bcrypt.GenerateFromPassword([]byte("paupau99"), bcrypt.DefaultCost)

// 		expectedData := user.Core{
// 			Email:       "mfauzanptra@gmail.com",
// 			Name:        "Muhamad Fauzan Putra",
// 			PhoneNumber: "085659171799",
// 			Password:    string(hashed),
// 			Address:     "Jln. Lembayung No 24, Bantul, Yogyakarta",
// 		}
// 		data.On("Login", email).Return(expectedData, nil)
// 		srv := New(data)
// 		token, res, err := srv.Login(email, password)
// 		assert.Nil(t, err)
// 		assert.Equal(t, expectedData.Email, res.Email)
// 		assert.NotNil(t, token)
// 		data.AssertExpectations(t)
// 	})

// 	t.Run("not found", func(t *testing.T) {
// 		inputEmail := "mfp@gmail.com"
// 		data.On("Login", inputEmail).Return(user.Core{}, errors.New("data not found"))

// 		srv := New(data)
// 		token, res, err := srv.Login(inputEmail, "be1422")
// 		assert.NotNil(t, err)
// 		assert.ErrorContains(t, err, "not found")
// 		assert.Empty(t, token)
// 		assert.Equal(t, uint(0), res.ID)
// 		data.AssertExpectations(t)
// 	})

// 	t.Run("wrong password", func(t *testing.T) {
// 		inputEmail := "mfauzanptra@gmail.com"
// 		hashed, _ := bcrypt.GenerateFromPassword([]byte("paupau99"), bcrypt.DefaultCost)
// 		expData := user.Core{
// 			Email:       "mfauzanptra@gmail.com",
// 			Name:        "Muhamad Fauzan Putra",
// 			PhoneNumber: "085659171799",
// 			Password:    string(hashed),
// 			Address:     "Jln. Lembayung No 24, Bantul, Yogyakarta",
// 		}
// 		data.On("Login", inputEmail).Return(expData, nil)

// 		srv := New(data)
// 		token, res, err := srv.Login(inputEmail, "be1423")
// 		assert.NotNil(t, err)
// 		assert.ErrorContains(t, err, "wrong password")
// 		assert.Empty(t, token)
// 		assert.Equal(t, uint(0), res.ID)
// 		data.AssertExpectations(t)
// 	})
// }

// func TestProfile(t *testing.T) {
// 	data := mocks.NewUserData(t)

// 	t.Run("Sukses lihat profile", func(t *testing.T) {
// 		resData := user.Core{ID: uint(1), Name: "jerry", Email: "jerry@alterra.id", PhoneNumber: "08123456"}

// 		data.On("Profile", uint(1)).Return(resData, nil).Once()

// 		srv := New(data)

// 		claims := jwt.MapClaims{}
// 		claims["authorized"] = true
// 		claims["userID"] = 1
// 		claims["exp"] = time.Now().Add(time.Hour * 1).Unix() //Token expires after 1 hour
// 		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// 		token.Valid = true

// 		res, err := srv.Profile(token)
// 		assert.Nil(t, err)
// 		assert.Equal(t, resData.ID, res.ID)
// 		data.AssertExpectations(t)
// 	})

// 	t.Run("jwt tidak valid", func(t *testing.T) {
// 		srv := New(data)

// 		_, token := helper.GenerateJWT(1)

// 		res, err := srv.Profile(token)
// 		assert.NotNil(t, err)
// 		assert.ErrorContains(t, err, "not found")
// 		assert.Equal(t, uint(0), res.ID)
// 	})

// 	t.Run("data tidak ditemukan", func(t *testing.T) {
// 		data.On("Profile", uint(4)).Return(user.Core{}, errors.New("data not found")).Once()

// 		srv := New(data)

// 		_, token := helper.GenerateJWT(4)
// 		pToken := token.(*jwt.Token)
// 		pToken.Valid = true
// 		res, err := srv.Profile(pToken)
// 		assert.NotNil(t, err)
// 		assert.ErrorContains(t, err, "not found")
// 		assert.Equal(t, uint(0), res.ID)
// 		data.AssertExpectations(t)
// 	})

// 	t.Run("masalah di server", func(t *testing.T) {
// 		data.On("Profile", mock.Anything).Return(user.Core{}, errors.New("terdapat masalah pada server")).Once()
// 		srv := New(data)

// 		_, token := helper.GenerateJWT(1)
// 		pToken := token.(*jwt.Token)
// 		pToken.Valid = true
// 		res, err := srv.Profile(pToken)
// 		assert.NotNil(t, err)
// 		assert.ErrorContains(t, err, "server")
// 		assert.Equal(t, uint(0), res.ID)
// 		data.AssertExpectations(t)
// 	})
// }

// func TestUpdate(t *testing.T) {
// 	data := mocks.NewUserData(t)
// 	updUser := user.Core{
// 		Name:        "fauzan",
// 		Email:       "mfauzanptra@gmail.com",
// 		PhoneNumber: "085659171799",
// 	}
// 	expectedData := user.Core{
// 		Name:        "fauzan",
// 		Email:       "mfauzanptra@gmail.com",
// 		PhoneNumber: "085659171799",
// 	}
// 	var a *multipart.FileHeader
// 	t.Run("update success", func(t *testing.T) {
// 		data.On("Profile", uint(1)).Return(expectedData, nil).Once()
// 		data.On("Update", uint(1), mock.Anything).Return(expectedData, nil).Once()
// 		srv := New(data)

// 		_, token := helper.GenerateJWT(1)
// 		pToken := token.(*jwt.Token)
// 		pToken.Valid = true

// 		res, err := srv.Update(pToken, updUser, a)
// 		assert.Nil(t, err)
// 		assert.Equal(t, expectedData.ID, res.ID)
// 		assert.Equal(t, expectedData.Name, res.Name)
// 		data.AssertExpectations(t)
// 	})

// 	t.Run("jwt not valid", func(t *testing.T) {
// 		srv := New(data)

// 		_, token := helper.GenerateJWT(1)

// 		res, err := srv.Update(token, updUser, a)
// 		assert.NotNil(t, err)
// 		assert.ErrorContains(t, err, "token error")
// 		assert.Equal(t, uint(0), res.ID)
// 	})

// 	t.Run("not found", func(t *testing.T) {
// 		data.On("Profile", uint(2)).Return(user.Core{}, errors.New("data not found")).Once()
// 		srv := New(data)

// 		_, token := helper.GenerateJWT(2)
// 		pToken := token.(*jwt.Token)
// 		pToken.Valid = true

// 		res, err := srv.Update(pToken, updUser, a)
// 		assert.NotNil(t, err)
// 		assert.Equal(t, res.Name, "")
// 		data.AssertExpectations(t)
// 	})

// 	t.Run("server problem", func(t *testing.T) {
// 		data.On("Profile", uint(1)).Return(user.Core{}, errors.New("server problem")).Once()
// 		srv := New(data)

// 		_, token := helper.GenerateJWT(1)
// 		pToken := token.(*jwt.Token)
// 		pToken.Valid = true

// 		res, err := srv.Update(pToken, updUser, a)
// 		assert.NotNil(t, err)
// 		assert.ErrorContains(t, err, "server")
// 		assert.Equal(t, res.Name, "")
// 		data.AssertExpectations(t)
// 	})

// 	t.Run("not found", func(t *testing.T) {
// 		data.On("Profile", uint(2)).Return(expectedData, nil).Once()
// 		data.On("Update", mock.Anything, mock.Anything).Return(user.Core{}, errors.New("data not found")).Once()
// 		srv := New(data)

// 		_, token := helper.GenerateJWT(2)
// 		pToken := token.(*jwt.Token)
// 		pToken.Valid = true

// 		res, err := srv.Update(pToken, updUser, a)
// 		assert.NotNil(t, err)
// 		assert.Equal(t, res.Name, "")
// 		data.AssertExpectations(t)
// 	})

// 	t.Run("server problem", func(t *testing.T) {
// 		data.On("Profile", uint(1)).Return(expectedData, nil).Once()
// 		data.On("Update", mock.Anything, mock.Anything).Return(user.Core{}, errors.New("server problem")).Once()
// 		srv := New(data)

// 		_, token := helper.GenerateJWT(1)
// 		pToken := token.(*jwt.Token)
// 		pToken.Valid = true

// 		res, err := srv.Update(pToken, updUser, a)
// 		assert.NotNil(t, err)
// 		assert.ErrorContains(t, err, "server")
// 		assert.Equal(t, res.Name, "")
// 		data.AssertExpectations(t)
// 	})
// }

// func TestDelete(t *testing.T) {
// 	data := mocks.NewUserData(t)
// 	t.Run("success delete", func(t *testing.T) {
// 		data.On("Delete", uint(1)).Return(nil).Once()

// 		srv := New(data)

// 		claims := jwt.MapClaims{}
// 		claims["authorized"] = true
// 		claims["userID"] = 1
// 		claims["exp"] = time.Now().Add(time.Hour * 1).Unix() //Token expires after 1 hour
// 		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// 		token.Valid = true

// 		err := srv.Delete(token)
// 		assert.Nil(t, err)
// 		data.AssertExpectations(t)
// 	})

// 	t.Run("jwt not valid", func(t *testing.T) {
// 		srv := New(data)

// 		_, token := helper.GenerateJWT(1)

// 		err := srv.Delete(token)
// 		assert.NotNil(t, err)
// 		assert.ErrorContains(t, err, "not found")
// 	})

// 	t.Run("data not found", func(t *testing.T) {
// 		data.On("Delete", uint(4)).Return(errors.New("data not found")).Once()

// 		srv := New(data)

// 		_, token := helper.GenerateJWT(4)
// 		pToken := token.(*jwt.Token)
// 		pToken.Valid = true
// 		err := srv.Delete(pToken)
// 		assert.NotNil(t, err)
// 		assert.ErrorContains(t, err, "not found")
// 		data.AssertExpectations(t)
// 	})

// 	t.Run("masalah di server", func(t *testing.T) {
// 		data.On("Delete", mock.Anything).Return(errors.New("terdapat masalah pada server")).Once()
// 		srv := New(data)

// 		_, token := helper.GenerateJWT(1)
// 		pToken := token.(*jwt.Token)
// 		pToken.Valid = true
// 		err := srv.Delete(pToken)
// 		assert.NotNil(t, err)
// 		assert.ErrorContains(t, err, "server")
// 		data.AssertExpectations(t)
// 	})
// }
