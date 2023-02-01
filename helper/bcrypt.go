package helper

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
)

func GeneratePassword(password string) string {
	var tmp1, tmp2 []byte
	hash := md5.Sum([]byte(password))
	num := len(hash) / 2
	tmp2, tmp1 = hash[:num], hash[len(hash)-num:]
	hash = md5.Sum([]byte(tmp1))
	tmp1 = hash[num:]
	hash = md5.Sum([]byte(tmp2))
	tmp2 = hash[:num]
	res := append(tmp2, tmp1...)
	return hex.EncodeToString(res[:])
}

func ComparePassword(hashedPassword, password string) error {
	password = GeneratePassword(password)
	if hashedPassword != password {
		return errors.New("password not match")
	}
	return nil
}
