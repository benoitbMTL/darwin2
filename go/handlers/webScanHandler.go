package handlers

import (
    "bytes"
    "os/exec"
    "github.com/labstack/echo/v4"
)

type RequestData struct {
    SelectedOption string `json:"selectedOption"`
}

func HandleWebScan(c echo.Context) error {
    var requestData RequestData
    if err := c.Bind(&requestData); err != nil {
        return echo.NewHTTPError(400, "Invalid data")
    }

    // Exécution de la commande nikto
    cmd := exec.Command("nikto", "-h", "https://192.168.4.10", "-timeout", "2", "-followredirects", "-until", "60s", "-useragent", "Nikto'\n''\r'X-Forwarded-For: 51.13.51.13")
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        return echo.NewHTTPError(500, "Error executing nikto command: "+err.Error())
    }

    // Renvoyer la sortie de la commande au client
    return c.String(200, out.String())
}
