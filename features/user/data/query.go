package data

import (
	"errors"
	"log"
	"sirloinapi/features/user"

	"gorm.io/gorm"
)

type userQry struct {
	db *gorm.DB
}

func New(db *gorm.DB) user.UserData {
	return &userQry{
		db: db,
	}
}

func (uq *userQry) Register(newUser user.Core) (user.Core, error) {
	cnv := CoreToData(newUser)
	err := uq.db.Create(&cnv).Error
	if err != nil {
		log.Println("error create query: ", err.Error())
		return user.Core{}, err
	}

	newUser.ID = cnv.ID

	return newUser, nil
}

func (uq *userQry) Login(email string) (user.Core, error) {
	res := User{}

	if err := uq.db.Where("email = ?", email).First(&res).Error; err != nil {
		log.Println("login query error: ", err.Error())
		return user.Core{}, errors.New("user not found")
	}
	return ToCore(res), nil
}

func (uq *userQry) Profile(id uint) (user.Core, error) {
	res := User{}
	if err := uq.db.Where("id = ?", id).First(&res).Error; err != nil {
		log.Println("Get By ID query error", err.Error())
		return user.Core{}, err
	}

	return ToCore(res), nil
}

func (uq *userQry) Update(id uint, updateData user.Core) (user.Core, error) {
	cnvUpd := CoreToData(updateData)
	qry := uq.db.Model(&User{}).Where("id = ?", id).Updates(cnvUpd)
	if err := qry.Error; err != nil {
		log.Println("error update user query : ", err)
		return updateData, err
	}
	return updateData, nil
}

func (uq *userQry) Delete(id uint) error {
	user := User{}
	qry := uq.db.Where("id = ?", id).Delete(&user)
	if affrows := qry.RowsAffected; affrows <= 0 {
		return errors.New("user doesn't exist")
	}
	if err := qry.Error; err != nil {
		log.Println(err)
		return err
	}
	return nil
}
