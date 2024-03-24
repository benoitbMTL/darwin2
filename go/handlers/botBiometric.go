package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleBiometric(c echo.Context) error {

	return c.String(http.StatusOK, "Bonjour")
}
