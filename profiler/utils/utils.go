package utils

import "log"

func LogError(err error, msg string) {
	if err != nil {
		log.Printf("\n%s: %s", msg, err)
	}
}
