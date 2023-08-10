package handler

import (
	"net/http"
	"net/mail"

	"github.com/labstack/echo/v4"
)

type testRequest struct {
	MonthlyIncome int `json:"monthly_income"`
	WorkingHour int `json:"working_hour"`
	Skill string `json:"skill"`
	RiskDegree int `json:"risk_degree"`
	IsInvestment bool `json:"is_investment"`
}

type blob []byte

type applyRequest struct {
	Name string `json:"name"`
	Mail mail.Address `json:"mail"`
	Artifact blob `json:"artifact"`
}

func Init(e *echo.Echo) {
	e.GET("/test", Test)
	e.GET("/calculate", Calculate)
	e.GET("/share", Share)
	e.GET("/apply", Apply)
}

func Test(c echo.Context) error {
	return nil
}

func Calculate(c echo.Context) error {
	var req testRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, req)
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