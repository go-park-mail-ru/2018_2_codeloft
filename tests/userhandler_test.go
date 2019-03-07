package tests
//
//import (
//	"bytes"
//	"encoding/json"
//	"github.com/go-park-mail-ru/2018_2_codeloft/handlers"
//	"github.com/go-park-mail-ru/2018_2_codeloft/models"
//	"net/http/httptest"
//	"testing"
//)
//
//var testUser = models.HelpUser{Login: "test", Password:"password"}
//
//
//func TestCreateUserOK(t *testing.T) {
//
//	testBody, err := json.Marshal(testUser)
//	testBuffer := bytes.NewBuffer(testBody)
//	if err != nil {
//		t.Fatal(err) // This should never happened if test setup correct
//	}
//	req := httptest.NewRequest("POST", "/user", testBuffer)
//	w := httptest.NewRecorder()
//	handler := handlers.UserHandler{}
//}