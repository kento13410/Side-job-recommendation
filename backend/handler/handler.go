package handler

import (
	"net/http"
	"net/mail"

	"github.com/labstack/echo/v4"
)

type testRequest struct {
	MonthlyIncome int `json:"monthly_income"`
	MonthlyWorkHours int `json:"monthly_work_hours"`
	Skill string `json:"skill"`
	SkillLevel int `json:"skill_level"`
	LeaningHours int `json:"leaning_hours"`
}

type blob []byte

type applyRequest struct {
	Name string `json:"name"`
	Mail mail.Address `json:"mail"`
	Artifact blob `json:"artifact"`
}

func Init(e *echo.Echo) {
	e.GET("/", Root)
	e.GET("/test", Test)
	e.POST("/calculate", Calculate)
	e.GET("/share", Share)
	e.POST("/apply", Apply)
}

func Root(c echo.Context) error {
	return c.File("../frontend/templates/index.html")
}

func Test(c echo.Context) error {
	return c.File("../frontend/templates/test.html")
}

func Calculate(c echo.Context) error {
	var req testRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.File("../frontend/templates/result.html")
}

func Share(c echo.Context) error {
	return nil
}

func Apply(c echo.Context) error {
	var req applyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, req)
}