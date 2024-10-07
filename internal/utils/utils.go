package utils

import (
	"fmt"
	"os"
)

func CheckError(err error, message string) {
	if err != nil {
		fmt.Printf("Error: %s - %v\n", message, err)
		os.Exit(1)
	}
}
