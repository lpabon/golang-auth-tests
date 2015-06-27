package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

/*
 * Works with https://github.com/auth0/go-jwt-middleware for Negroni
 */

func main() {
	token := jwt.New(jwt.SigningMethodHS256)
	s, err := token.SignedString([]byte("My Secret"))
	if err != nil {
		fmt.Print(err)
	} else {
		fmt.Println(s)
	}
}
