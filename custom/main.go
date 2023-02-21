package main

import (
	"net/http"

	"github.com/dewidyabagus/echo-validator/custom/validator"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type User struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"email"`
	Role  int    `json:"role" validate:"required,gte=1,lte=10"`
}

const address = ":7001"

func main() {
	// Instance engine validator dengan menggunakan custom tag name function
	validation := validator.New(validator.Options{TagNameFunc: "json"})

	// Instance HTTP service
	e := echo.New()
	e.Validator = validation
	e.HTTPErrorHandler = CustomHTTPErrorWithValTranslator(validation.ErrorFormTranslator)
	e.Use(middleware.Recover()) // Recovery panic error dengan default hanya logging saja

	// HTTP Handlers
	e.POST("/users", func(c echo.Context) error {
		user := new(User)
		if err := c.Bind(user); err != nil {
			return err
		}
		if err := c.Validate(user); err != nil {
			return err
		}
		return c.JSON(http.StatusCreated, echo.Map{"message": "berhasil membuat user"})
	})

	e.Start(address)
}
