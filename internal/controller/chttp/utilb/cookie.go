package utilb

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

func DeleteCookie(w http.ResponseWriter, cookie *http.Cookie) {
	cookie.Expires = time.Unix(0, 0)
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}

func RandomCookieValue() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf), nil
}
