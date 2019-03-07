package tests
// Maintainer D. Anokhin
import (
	"encoding/hex"
	"github.com/go-park-mail-ru/2018_2_codeloft/services"
	"testing"
	"time"
)

func TestCookieHasHexValue(t *testing.T) {
	cookie := services.GenerateCookie("anything")
	_, err := hex.DecodeString(cookie.Value)
	if err != nil {
		t.Errorf("Cookie value: \"%s\" is not hex", cookie.Value)
	}
}

func TestCookieHasCurrentExpireDate(t *testing.T) {
	expiresDate := time.Now().Add(30 * 24 * time.Hour)
	cookie := services.GenerateCookie("anything")
	if cookie.Expires.Truncate(time.Second) != expiresDate.Truncate(time.Second) {
		t.Errorf("Expected expires: %v, cookie expires: %v", expiresDate, cookie.Expires)
	}
}

