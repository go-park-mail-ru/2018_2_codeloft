package tests

import (
	"testing"

	"github.com/go-park-mail-ru/2018_2_codeloft/validator"
)

func TestEmailOk(t *testing.T) {
	in := "123@ya.ru"
	if err := validator.ValidateEmail(in); err != nil {
		t.Error("Unexpected error ", err)
	}
}

func TestEmailLong(t *testing.T) {
	in := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	if err := validator.ValidateEmail(in); err.Error() != validator.ErrorLong.Error() {
		t.Errorf("Got %v\n Expected %v", err, validator.ErrorLong)
	}
}

func TestPasswordLong(t *testing.T) {
	in := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	if err := validator.ValidatePassword(in); err == nil || err.Error() != validator.ErrorLong.Error() {
		t.Errorf("Got %v\n Expected %v", err, validator.ErrorLong)
	}
}

func TestPasswordShort(t *testing.T) {
	in := "aaaaaa"
	if err := validator.ValidatePassword(in); err == nil || err.Error() != validator.ErrorShort.Error() {
		t.Errorf("Got %v\n Expected %v", err, validator.ErrorShort)
	}
}

func TestPasswordInvalidChar(t *testing.T) {
	in := "a8fjblddzi_a"
	if err := validator.ValidatePassword(in); err == nil || err.Error() != validator.ErrorInvalidChar.Error() {
		t.Errorf("Got %v\n Expected %v", err, validator.ErrorShort)
	}
}

func TestPasswordOk(t *testing.T) {
	in := "8abCd1eZ"
	if err := validator.ValidatePassword(in); err != nil {
		t.Error(err)
	}
}

func TestLoginLong(t *testing.T) {
	in := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	if err := validator.ValidateLogin(in); err == nil || err.Error() != validator.ErrorLong.Error() {
		t.Errorf("Got %v\n Expected %v", err, validator.ErrorLong)
	}
}

func TestLoginInvalidChar(t *testing.T) {
	in := "Bak;"
	if err := validator.ValidateLogin(in); err == nil || err.Error() != validator.ErrorInvalidChar.Error() {
		t.Errorf("Got %v\n Expected %v", err, validator.ErrorInvalidChar)
	}
}

func TestLoginOk(t *testing.T) {
	in := "Alex_Bak.2"
	if err := validator.ValidateLogin(in); err != nil {
		t.Error(err)
	}
}
