package utils

import (
	"math/rand"
	"os"
)

var alphabet = "2345678923456789234567892345678923456789qwertyuipasdfghjkxcvbnm"
var idLength = 10

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func RandomID() string {
	ID := make([]byte, 0)
	abLen := len(alphabet)
	for i := 0; i < idLength; i++ {
		ID = append(ID, alphabet[rand.Intn(abLen)])
	}
	return string(ID)
}
