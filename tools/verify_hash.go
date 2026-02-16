package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	hash := "$2a$10$jnJdE8SwgG1r7CqApYJYiu7Ti9Pz8q/b2oLdVaF2r/rE7hIqVwBXO"
	pass := "123456"
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err == nil {
		fmt.Println("Hash matches 123456")
	} else {
		fmt.Println("Hash DOES NOT match 123456:", err)
	}
}
