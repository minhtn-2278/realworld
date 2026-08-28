package handlers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

func newHandlerTestContext(method string, target string, body string) (*echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	request := httptest.NewRequest(method, target, strings.NewReader(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()

	return e.NewContext(request, recorder), recorder
}

func setPathParams(context *echo.Context, params ...string) {
	if len(params) == 0 {
		return
	}

	pathValues := make(echo.PathValues, 0, len(params)/2)
	for index := 0; index < len(params); index += 2 {
		pathValues = append(pathValues, echo.PathValue{Name: params[index], Value: params[index+1]})
	}
	context.SetPathValues(pathValues)
}

func setAuthenticatedUser(context *echo.Context) {
	context.Set("user", &jwt.Token{Claims: jwt.MapClaims{
		"sub":        "1",
		"username":   "alice",
		"token_type": "access",
	}})
}

func assertHandlerStatus(t *testing.T, err error, recorder *httptest.ResponseRecorder, expectedStatus int) {
	t.Helper()
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if recorder.Code != expectedStatus {
		t.Fatalf("status code = %d, want %d; body = %s", recorder.Code, expectedStatus, recorder.Body.String())
	}
}
