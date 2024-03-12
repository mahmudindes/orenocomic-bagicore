package utila

import "math/rand"

var (
	RandomStringNumber    = "0123456789"
	RandomStringLowercase = "abcdefghijklmnopqrstuvwxyz"
	RandomStringUppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	RandomStringGeneral   = RandomStringNumber + RandomStringLowercase + RandomStringUppercase
)

func RandomString(charset string, length int) string {
	buf := make([]byte, length)
	for i := range buf {
		buf[i] = charset[rand.Intn(len(charset))]
	}
	return string(buf)
}
