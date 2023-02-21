package main

import (
	"net/http"
	"time"

	validator "github.com/go-playground/validator/v10"
	echo "github.com/labstack/echo/v4"
)

type response struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Validation interface{} `json:"validation,omitempty"`
	Timestamp  int64       `json:"timestamp"`
}

type validatorFunc func(valErrors validator.ValidationErrors) interface{}

// Membuat custom error dengan tambahan untuk jenis error dari proses validasi data
// akan diproses terpisah sesuai dengan aturan pada parameter function
func CustomHTTPErrorWithValTranslator(validFunc validatorFunc) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		var res response
		var report *echo.HTTPError

		switch val := err.(type) {
		default:
			report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())

		case *echo.HTTPError:
			report = val

		case validator.ValidationErrors:
			res.Validation = validFunc(val)
			report = echo.NewHTTPError(http.StatusBadRequest, "invalid data")
		}

		res.Code = report.Code
		res.Message, _ = report.Message.(string)
		res.Timestamp = time.Now().UnixMilli()

		c.JSON(res.Code, res)
	}
}
