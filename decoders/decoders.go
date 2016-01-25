package decoders

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo"
)

const minNameLength = 3
const minPasswordLength = 6

// store the errors as a map
type Errors map[string]string

func (e Errors) Error() string {
	return fmt.Sprintf("%v", e)
}

func (e Errors) Render(c *echo.Context) error {
	return c.JSON(http.StatusBadRequest, e)
}

type Decoder interface {
	Decode() error
}

// Decode decodes JSON body of request and runs through validator
// if err, ok := d.Decode(c); !ok {
// return err }

func Decode(c *echo.Context, decoder Decoder) (error, bool) {
	if err := c.Bind(decoder); err != nil {
		return err, false
	}
	if err := decoder.Decode(); err != nil {
		if e, ok := err.(Errors); ok {
			return e.Render(c), false
		}
		return err, false
	}
	return nil, true
}

type RecoverPassword struct {
	Identifier string `json:"identifier"`
}

func (r *RecoverPassword) Decode() error {

	if r.Identifier == "" {
		return Errors{
			"identifier": "Email or user name required",
		}
	}
	return nil
}

type Signup struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Signup) Decode() error {

	s.Name = strings.Trim(s.Name, " ")
	s.Email = strings.ToLower(strings.Trim(s.Email, " "))

	errors := make(Errors)

	if len(s.Name) < minNameLength {
		errors["name"] = fmt.Sprintf("Your name must be at least %d characters long", minNameLength)
	}
	if !govalidator.IsEmail(s.Email) {
		errors["email"] = "Email address is required"
	}
	if len(s.Password) < minPasswordLength {
		errors["password"] = fmt.Sprintf("Password must be at least %d characters long", minPasswordLength)
	}
	if len(errors) > 0 {
		return errors
	}

	return nil
}

type Login struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

func (l *Login) Decode() error {
	errors := make(Errors)
	if l.Identifier == "" {
		errors["identifier"] = "Email or user name missing"
	}
	if l.Password == "" {
		errors["password"] = "Password missing"
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}

type NewChannel struct {
	URL string `json:"url"`
}

func (n *NewChannel) Decode() error {
	n.URL = strings.Trim(n.URL, " ")
	if n.URL == "" || !govalidator.IsURL(n.URL) {
		return Errors{
			"url": "Valid URL is required",
		}
	}
	return nil
}

type NewEmail struct {
	Email string `json:"email"`
}

func (n *NewEmail) Decode() error {
	n.Email = strings.Trim(strings.ToLower(n.Email), " ")
	if n.Email == "" || !govalidator.IsEmail(n.Email) {
		return Errors{
			"email": "Valid email is required",
		}
	}
	return nil
}

type NewPassword struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (n *NewPassword) Decode() error {
	errors := make(Errors)
	if n.OldPassword == "" {
		errors["oldPassword"] = "Old password is required"
	}
	if len(n.NewPassword) < minPasswordLength {
		errors["newPassword"] = fmt.Sprintf("New password must be at least %d characters long", minPasswordLength)
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}
