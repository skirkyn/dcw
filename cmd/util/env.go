package util

import (
	"log"
	"os"
	"strconv"
)

func GetEnvString(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
func GetEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		intVal, err := strconv.Atoi(value)
		if err != nil {
			log.Println(err)
			return fallback
		}
		return intVal
	}
	return fallback
}
