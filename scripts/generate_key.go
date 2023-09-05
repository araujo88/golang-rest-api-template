package main

import (
	"fmt"
	"golang-rest-api-template/pkg/auth"
)

func main() {
	fmt.Println(auth.GenerateRandomKey())
}
