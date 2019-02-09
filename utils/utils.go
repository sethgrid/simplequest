package utils

import "log"

func StrInList(needle string, haystack []string) bool {
	for _, element := range haystack {
		if needle == element {
			return true
		}
	}
	return false
}

var Debug = true

func Debugf(format string, v ...interface{}) {
	if !Debug {
		return
	}
	log.Printf(format, v...)
}
