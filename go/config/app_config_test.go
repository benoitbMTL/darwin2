package config

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestSaveAndReadProfileWithoutActivatingIt(t *testing.T) {
	Initialize()
	e := echo.New()

	body := `{"NAME":"Saved profile","DVWAURL":"https://saved.example.test"}`
	request := httptest.NewRequest(http.MethodPost, "/save-config", strings.NewReader(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	if err := SaveConfig(e.NewContext(request, recorder)); err != nil {
		t.Fatalf("SaveConfig returned an error: %v", err)
	}

	if active := GetCurrentConfig().NAME; active != "Default" {
		t.Fatalf("active profile = %q, want Default", active)
	}

	request = httptest.NewRequest(http.MethodGet, "/configs/Saved%20profile", nil)
	recorder = httptest.NewRecorder()
	context := e.NewContext(request, recorder)
	context.SetPath("/configs/:name")
	context.SetParamNames("name")
	context.SetParamValues("Saved profile")
	if err := GetConfigByName(context); err != nil {
		t.Fatalf("GetConfigByName returned an error: %v", err)
	}
	if !strings.Contains(recorder.Body.String(), `"NAME":"Saved profile"`) {
		t.Fatalf("unexpected profile response: %s", recorder.Body.String())
	}
}
