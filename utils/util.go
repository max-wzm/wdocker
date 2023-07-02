package utils

import (
	"math/rand"
	"os"
	"os/exec"
)

var alphabet = "2345678923456789234567892345678923456789qwertyuipasdfghjkxcvbnm"
var idLength = 10

func RandomID() string {
	ID := make([]byte, 0)
	abLen := len(alphabet)
	for i := 0; i < idLength; i++ {
		ID = append(ID, alphabet[rand.Intn(abLen)])
	}
	return string(ID)
}

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

func CmdRunStd(name string, arg ...string) error {
	c := exec.Command(name, arg...)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	err := c.Run()
	return err
}
