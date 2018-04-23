package util

import (
	"log"
	"runtime"
)

func CheckErr(err error) bool {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.Println(file, line, ":", err)
		return true
	}
	return false
}
