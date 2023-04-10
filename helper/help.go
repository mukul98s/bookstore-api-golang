package helper

import "fmt"

func CheckError(err error, message string) {
	if err != nil {
		fmt.Println("Error:", message)
		panic(err)
	}
}
