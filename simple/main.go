package main

import (
	"fmt"
	"net/http"

	validator "github.com/go-playground/validator/v10"
	echo "github.com/labstack/echo/v4"
)

type User struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type Validation struct {
	validator *validator.Validate
}

func (v *Validation) Validate(s interface{}) error {
	return v.validator.Struct(s)
}

func main() {
	e := echo.New()
	e.Validator = &Validation{validator: validator.New()}

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		var report *echo.HTTPError

		switch valErr := err.(type) {
		default:
			report = echo.NewHTTPError(http.StatusInternalServerError, valErr.Error())

		case validator.ValidationErrors:
			report = echo.NewHTTPError(http.StatusBadRequest)
			if len(valErr) > 0 {
				switch valErr[0].Tag() {
				case "required":
					report.Message = fmt.Sprintf("%s is required", valErr[0].Field())

				case "email":
					report.Message = fmt.Sprintf("%s is not valid email", valErr[0].Field())

				default:
					report.Message = valErr[0].Error()
				}
			} else {
				report.Message = valErr.Error()
			}

		case *echo.HTTPError:
			report = valErr
		}

		c.JSON(report.Code, report)
	}

	e.POST("/users", func(c echo.Context) error {
		user := new(User)
		if err := c.Bind(user); err != nil {
			return err
		}
		if err := c.Validate(user); err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, echo.Map{"message": "berhasil membuat user baru!"})
	})

	e.Start(":7001")
}
