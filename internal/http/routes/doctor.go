package routes

import (
	"net/http"

	"github.com/ambientis-org/hefesto/internal/db/postgres/models"
	"github.com/labstack/echo/v4"
)

func registerDoctor(c echo.Context) error {
	doctor := &models.Doctor{}
	err := c.Bind(doctor)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	DataBase.Create(doctor)

	return c.JSON(http.StatusCreated, doctor)
}

func getDoctor(c echo.Context) error {
	d := &models.Doctor{}
	DataBase.Where("id = ?", c.QueryParam("doctor_id")).First(d)

	return c.JSON(http.StatusOK, d)
}

// setupLogin Add Login handlers to API
func (router *Router) setupDoctor() {
	router.addGroup("/doctor")

	doctorGroup := API.groups["/doctor"]
	doctorGroup.POST("", registerDoctor)
	doctorGroup.GET("", getDoctor)
}
