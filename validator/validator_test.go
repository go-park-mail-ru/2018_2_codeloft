package validator

import (
	"testing"
)

func TestEmailOk(t *testing.T) {
	in := "123@ya.ru"
	if err := ValidateEmail(in); err != nil {
		t.Error("Unexpected error ", err)
	}
}

func TestEmailLong(t *testing.T) {
	in := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	if err := ValidateEmail(in); err.Error() != errorLong.Error() {
		t.Errorf("Got %v\n Expected %v", err, errorLong)
	}
}

func TestPasswordLong(t *testing.T) {
	in := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	if err := ValidatePassword(in); err == nil || err.Error() != errorLong.Error() {
		t.Errorf("Got %v\n Expected %v", err, errorLong)
	}
}

func TestPasswordShort(t *testing.T) {
	in := "aaaaaa"
	if err := ValidatePassword(in); err == nil || err.Error() != errorShort.Error() {
		t.Errorf("Got %v\n Expected %v", err, errorShort)
	}
}

func TestPasswordInvalidChar(t *testing.T) {
	in := "a8fjblddzi_a"
	if err := ValidatePassword(in); err == nil || err.Error() != errorInvalidChar.Error() {
		t.Errorf("Got %v\n Expected %v", err, errorShort)
	}
}

func TestPasswordOk(t *testing.T) {
	in := "8abCd1eZ"
	if err := ValidatePassword(in); err != nil {
		t.Error(err)
	}
}

func TestLoginLong(t *testing.T) {
	in := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	if err := ValidateLogin(in); err == nil || err.Error() != errorLong.Error() {
		t.Errorf("Got %v\n Expected %v", err, errorLong)
	}
}

func TestLoginInvalidChar(t *testing.T) {
	in := "Bak;"
	if err := ValidateLogin(in); err == nil || err.Error() != errorInvalidChar.Error() {
		t.Errorf("Got %v\n Expected %v", err, errorInvalidChar)
	}
}

func TestLoginOk(t *testing.T) {
	in := "Alex_Bak.2"
	if err := ValidateLogin(in); err != nil {
		t.Error(err)
	}
}
