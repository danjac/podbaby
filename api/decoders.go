package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/asaskevich/govalidator"
	"github.com/danjac/podbaby/api/Godeps/_workspace/src/github.com/labstack/echo"
)

const (
	minNameLength     = 3
	minPasswordLength = 6
)

type validator struct {
	context *echo.Context
	errors  map[string]string
}

func newValidator(c *echo.Context) *validator {
	errors := make(map[string]string)
	return &validator{c, errors}
}

func (v *validator) invalid(field, msg string) *validator {
	v.errors[field] = msg
	return v
}

func (v *validator) render() error {
	return v.context.JSON(http.StatusBadRequest, v.errors)
}

func (v *validator) ok() bool {
	return len(v.errors) == 0
}

func (v *validator) validate(d decoder) (bool, error) {
	if err := v.context.Bind(d); err != nil {
		return false, err
	}
	if d.decode(v); !v.ok() {
		return false, v.render()
	}
	return true, nil
}

type decoder interface {
	decode(*validator)
}

type recoverPasswordDecoder struct {
	Identifier string `json:"identifier"`
}

func (d *recoverPasswordDecoder) decode(v *validator) {
	if d.Identifier == "" {
		v.invalid("identifier", "Email or user name required")
	}
}

type signupDecoder struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (d *signupDecoder) decode(v *validator) {

	d.Name = strings.Trim(d.Name, " ")
	d.Email = strings.ToLower(strings.Trim(d.Email, " "))

	if len(d.Name) < minNameLength {
		v.invalid("name", fmt.Sprintf("Your name must be at least %d characters long", minNameLength))
	}
	if !govalidator.IsEmail(d.Email) {
		v.invalid("email", "Email address is required")
	}
	if len(d.Password) < minPasswordLength {
		v.invalid("password", fmt.Sprintf("Password must be at least %d characters long", minPasswordLength))
	}
}

type loginDecoder struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

func (d *loginDecoder) decode(v *validator) {
	if d.Identifier == "" {
		v.invalid("identifier", "Email or user name missing")
	}
	if d.Password == "" {
		v.invalid("password", "Password missing")
	}
}

type newChannelDecoder struct {
	URL string `json:"url"`
}

func (d *newChannelDecoder) decode(v *validator) {
	d.URL = strings.Trim(d.URL, " ")
	if d.URL == "" || !govalidator.IsURL(d.URL) {
		v.invalid("url", "Valid URL is required")
	}
}

type changeEmailDecoder struct {
	Email string `json:"email"`
}

func (d *changeEmailDecoder) decode(v *validator) {
	d.Email = strings.Trim(strings.ToLower(d.Email), " ")
	if d.Email == "" || !govalidator.IsEmail(d.Email) {
		v.invalid("email", "Valid email is required")
	}
}

type changePasswordDecoder struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (d *changePasswordDecoder) decode(v *validator) {
	if d.OldPassword == "" {
		v.invalid("oldPassword", "Old password is required")
	}
	if len(d.NewPassword) < minPasswordLength {
		v.invalid("newPassword", fmt.Sprintf("New password must be at least %d characters long", minPasswordLength))
	}
}
