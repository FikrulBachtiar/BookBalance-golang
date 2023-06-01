package utils

import "strings"

func TrimPhoneNumber(number string) string {
	return strings.Replace(number, "+", "", -1);
}