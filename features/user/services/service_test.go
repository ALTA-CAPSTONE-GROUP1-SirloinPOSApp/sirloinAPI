package service

import (
	"errors"
	"sirloinapi/features/user"
	"sirloinapi/helper"
	"sirloinapi/mocks"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestLogn(t *testing.T) {
	data := mocks.NewUserData(t)
	// input dan respond untuk mock data
	password := "Amr12345"
	hashed := helper.GeneratePassword(password)
	inputData := user.Core{
		Email:    "jerr@alterra.id",
		Password: password,
	}
	expectedData := user.Core{
		Email:        "mfauzanptra@gmail.com",
		BusinessName: "Muhamad Fauzan Putra",
		PhoneNumber:  "085659171799",
		Password:     hashed,
		Address:      "Jln. Lembayung No 24, Bantul, Yogyakarta",
	}
	t.Run("succcess login", func(t *testing.T) {

		// res dari data akan mengembalik password yang sudah di hash
		data.On("Login", inputData.Email).Return(expectedData, nil).Once()
		srv := New(data)
		inputData.Password = password
		token, res, err := srv.Login(inputData.Email, inputData.Password)
		assert.Nil(t, err)
		assert.Equal(t, expectedData.Email, res.Email)
		assert.NotNil(t, token)
		data.AssertExpectations(t)
	})

	t.Run("server problem", func(t *testing.T) {
		data.On("Login", inputData.Email).Return(user.Core{}, errors.New("server error")).Once()
		srv := New(data)
		token, res, err := srv.Login(inputData.Email, inputData.Password)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Empty(t, token)
		assert.Equal(t, uint(0), res.ID)
		data.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		data.On("Login", inputData.Email).Return(user.Core{}, errors.New("data not found")).Once()
		srv := New(data)
		token, res, err := srv.Login(inputData.Email, inputData.Password)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
		assert.Empty(t, token)
		assert.Equal(t, uint(0), res.ID)
		data.AssertExpectations(t)
	})

	t.Run("wrong password", func(t *testing.T) {
		data.On("Login", inputData.Email).Return(expectedData, nil)
		srv := New(data)
		token, res, err := srv.Login(inputData.Email, "Abe12345")
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "wrong password")
		assert.Empty(t, token)
		assert.Equal(t, uint(0), res.ID)
		data.AssertExpectations(t)
	})

}

func TestProfile(t *testing.T) {
	data := mocks.NewUserData(t)
	expectedData := user.Core{
		Email:        "mfauzanptra@gmail.com",
		BusinessName: "Muhamad Fauzan Putra",
		PhoneNumber:  "085659171799",
		Address:      "Jln. Lembayung No 24, Bantul, Yogyakarta",
	}
	t.Run("Sukses lihat profile", func(t *testing.T) {
		data.On("Profile", uint(1)).Return(expectedData, nil).Once()
		srv := New(data)
		claims := jwt.MapClaims{}
		claims["authorized"] = true
		claims["userID"] = 1
		claims["exp"] = time.Now().Add(time.Hour * 1).Unix() //Token expires after 1 hour
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		token.Valid = true

		res, err := srv.Profile(token)
		assert.Nil(t, err)
		assert.Equal(t, expectedData.ID, res.ID)
		data.AssertExpectations(t)
	})

	t.Run("jwt tidak valid", func(t *testing.T) {
		srv := New(data)

		_, token := helper.GenerateJWT(1)

		res, err := srv.Profile(token)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
		assert.Equal(t, uint(0), res.ID)
	})

	t.Run("data tidak ditemukan", func(t *testing.T) {
		data.On("Profile", uint(4)).Return(user.Core{}, errors.New("data not found")).Once()

		srv := New(data)

		_, token := helper.GenerateJWT(4)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Profile(pToken)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
		assert.Equal(t, uint(0), res.ID)
		data.AssertExpectations(t)
	})

	t.Run("masalah di server", func(t *testing.T) {
		data.On("Profile", mock.Anything).Return(user.Core{}, errors.New("terdapat masalah pada server")).Once()
		srv := New(data)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Profile(pToken)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Equal(t, uint(0), res.ID)
		data.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	data := mocks.NewUserData(t)
	password := "Amr12345"
	hash := helper.GeneratePassword(password)
	updUser := user.Core{
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

	t.Run("not found", func(t *testing.T) {
		data.On("Update", uint(1), updUser).Return(user.Core{}, errors.New("data not found")).Once()
		srv := New(data)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.Update(pToken, updUser)
		assert.NotNil(t, err)
		assert.Equal(t, res.BusinessName, "")
		data.AssertExpectations(t)
	})

	t.Run("server problem", func(t *testing.T) {
		data.On("Update", uint(1), updUser).Return(user.Core{}, errors.New("server problem")).Once()
		srv := New(data)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.Update(pToken, updUser)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Equal(t, res.BusinessName, "")
		data.AssertExpectations(t)
	})

	t.Run("update success", func(t *testing.T) {
		updUser.Password = hash
		data.On("Update", uint(1), updUser).Return(expectedData, nil).Once()
		srv := New(data)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		updUser.Password = password
		res, err := srv.Update(pToken, updUser)
		assert.Nil(t, err)
		assert.Equal(t, expectedData.ID, res.ID)
		assert.Equal(t, expectedData.BusinessName, res.BusinessName)
		data.AssertExpectations(t)
	})

	t.Run("user / email already exist", func(t *testing.T) {
		updUser.Password = hash
		data.On("Update", uint(1), updUser).Return(user.Core{}, errors.New("Duplicate users.email")).Once()
		srv := New(data)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		updUser.Password = password
		res, err := srv.Update(pToken, updUser)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "user already exist")
		assert.Equal(t, res.BusinessName, "")
		data.AssertExpectations(t)
	})

	t.Run("phone number already exist", func(t *testing.T) {
		updUser.Password = hash
		data.On("Update", uint(1), updUser).Return(user.Core{}, errors.New("Duplicate users.phone_number")).Once()
		srv := New(data)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		updUser.Password = password
		res, err := srv.Update(pToken, updUser)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "phone number already exist")
		assert.Equal(t, res.BusinessName, "")
		data.AssertExpectations(t)
	})

	t.Run("Password validasi", func(t *testing.T) {
		srv := New(data)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		updUser.Password = "123"
		res, err := srv.Update(pToken, updUser)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "password")
		assert.Equal(t, res.BusinessName, "")
	})

	t.Run("Business name validasi", func(t *testing.T) {
		srv := New(data)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		updUser.Password = password
		updUser.BusinessName = "AMR#%^%"
		res, err := srv.Update(pToken, updUser)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "alpha_space")
		assert.Equal(t, res.BusinessName, "")
	})

	t.Run("Email validasi", func(t *testing.T) {
		srv := New(data)
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		updUser.Password = password
		updUser.BusinessName = "Toko Ira"
		updUser.Email = "tokoira.com"
		res, err := srv.Update(pToken, updUser)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "email")
		assert.Equal(t, res.BusinessName, "")
	})

	t.Run("Phone number validasi", func(t *testing.T) {
		srv := New(data)
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		updUser.Password = password
		updUser.BusinessName = "Toko Ira"
		updUser.Email = "tokoira@gmail.com"
		updUser.PhoneNumber = "oiajsdojhasodi"
		res, err := srv.Update(pToken, updUser)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "PhoneNumber")
		assert.Equal(t, res.BusinessName, "")
	})

	t.Run("jwt not valid", func(t *testing.T) {
		srv := New(data)

		_, token := helper.GenerateJWT(1)

		res, err := srv.Update(token, updUser)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "token error")
		assert.Equal(t, uint(0), res.ID)
	})

}

func TestDelete(t *testing.T) {
	data := mocks.NewUserData(t)
	t.Run("success delete", func(t *testing.T) {
		data.On("Delete", uint(1)).Return(nil).Once()

		srv := New(data)

		claims := jwt.MapClaims{}
		claims["authorized"] = true
		claims["userID"] = 1
		claims["exp"] = time.Now().Add(time.Hour * 1).Unix() //Token expires after 1 hour
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		token.Valid = true

		err := srv.Delete(token)
		assert.Nil(t, err)
		data.AssertExpectations(t)
	})

	t.Run("jwt not valid", func(t *testing.T) {
		srv := New(data)

		_, token := helper.GenerateJWT(1)

		err := srv.Delete(token)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
	})

	t.Run("data not found", func(t *testing.T) {
		data.On("Delete", uint(4)).Return(errors.New("data not found")).Once()

		srv := New(data)

		_, token := helper.GenerateJWT(4)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Delete(pToken)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
		data.AssertExpectations(t)
	})

	t.Run("masalah di server", func(t *testing.T) {
		data.On("Delete", mock.Anything).Return(errors.New("terdapat masalah pada server")).Once()
		srv := New(data)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Delete(pToken)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		data.AssertExpectations(t)
	})
}

func TestRegisterDevice(t *testing.T) {
	data := mocks.NewUserData(t)
	dvctoken := "udaiowbhnduiwiudhiaudwbiwafbiawfbfw"
	srv := New(data)

	t.Run("success register device", func(t *testing.T) {
		data.On("RegisterDevice", uint(1), dvctoken).Return(nil).Once()
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		err := srv.RegisterDevice(pToken, dvctoken)
		assert.Nil(t, err)
		data.AssertExpectations(t)
	})

	t.Run("jwt not valid", func(t *testing.T) {
		_, token := helper.GenerateJWT(1)

		err := srv.RegisterDevice(token, dvctoken)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
	})

	t.Run("error not found", func(t *testing.T) {
		data.On("RegisterDevice", uint(1), dvctoken).Return(errors.New("not found")).Once()
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		err := srv.RegisterDevice(pToken, dvctoken)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "found")
		data.AssertExpectations(t)
	})

	t.Run("error duplicated", func(t *testing.T) {
		data.On("RegisterDevice", uint(1), dvctoken).Return(errors.New("duplicated")).Once()
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		err := srv.RegisterDevice(pToken, dvctoken)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "duplicated")
		data.AssertExpectations(t)
	})

	t.Run("server problem", func(t *testing.T) {
		data.On("RegisterDevice", uint(1), dvctoken).Return(errors.New("server")).Once()
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		err := srv.RegisterDevice(pToken, dvctoken)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "server")
		data.AssertExpectations(t)
	})
}

func Test(t *testing.T) {
	data := mocks.NewUserData(t)
	srv := New(data)

	t.Run("success unreg device", func(t *testing.T) {
		data.On("UnregDevice", uint(1)).Return(nil).Once()
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		err := srv.UnregDevice(pToken)
		assert.Nil(t, err)
		data.AssertExpectations(t)
	})

	t.Run("jwt not valid", func(t *testing.T) {
		_, token := helper.GenerateJWT(1)

		err := srv.UnregDevice(token)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
	})

	t.Run("error not found", func(t *testing.T) {
		data.On("UnregDevice", uint(1)).Return(errors.New("not found")).Once()
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		err := srv.UnregDevice(pToken)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "found")
		data.AssertExpectations(t)
	})

	t.Run("server problem", func(t *testing.T) {
		data.On("UnregDevice", uint(1)).Return(errors.New("server")).Once()
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		err := srv.UnregDevice(pToken)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "server")
		data.AssertExpectations(t)
	})
}
