package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(providedPwd string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(providedPwd), 14)
	if err != nil {
		fmt.Println(err)
	}
	return string(bytes)
}
