# Custom Validate Data Echo Framework
Framework Echo tidak mendukung proses validasi data secara bawaan. Ketika kita ingin melakukan validasi data melalui `echo.Context` yang sudah disediakan framework echo maka harus melakukan registrasi (overwrite) method `Validator` dengan custom validator yang digunakan.

Sebagai contoh saya membuat sebuah endpoint yang digunakan untuk registrasi user dengan URL `POST http://localhost:7001/users`, dimana endpoint tersebut menerima data dalam format `JSON` dengan field `name` dan `email` user. Untuk masing-masing data yang diterima wajib diisi dan untuk `email` harus sesuai dengan format email. Untuk menulis rules menggunakan tag `validate:"rules"`.   
```go
type User struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}
```
Untuk lebih detail macam-macam validasi dan metode yang digunakan silahkan lihat dilaman github https://github.com/go-playground/validator.

Proses penulisan custom validator dimulai dengan membuat `struct` yang memiliki field bertipe `*validator.Validate` dan method `Validate(s interface{}) error` yang nantinya akan digunakan sebagai proses validasi data melalui `echo.Context`.
```go
type Validation struct {
	validator *validator.Validate
}

func (v *Validation) Validate(s interface{}) error {
	return v.validator.Struct(s)
}
```

Proses overwrite validator Echo dapat langsung dilakukan dengan menggunakan field instance Echo dan mengisinya dengan pointer struct `Validation` (untuk field `validator` pastikan sudah menerima instance dari engine validator):   
```go
func main() {
    e := echo.New()
    e.Validator = &Validation{validator: validator.New()}
}
```

Untuk custom response menggunakan `echo.HTTPErrorHandler` yang sudah dikelompokan berdasarkan jenis error yang bersumber dari framework Echo (HTTPError), proses validasi atau lainnya.
```go
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
```

Dari penjelasan diatas ketika disusun menjadi kode yang utuh seperti berikut :
```go
// github.com/dewidyabagus/echo-validator/simple/main.go
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
```

Menjalankan HTTP service:
```bash
go run ./simple
```

Hasil ketika data email tidak valid:
```bash
curl -X POST http://localhost:7001/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Joe","email":"joe@invalid-domain"}'
{"message":"Email is not valid email"}
```

Referensi :
- https://dasarpemrogramangolang.novalagung.com/C-http-request-payload-validation.html
- https://echo.labstack.com/guide/request/#validate-data
