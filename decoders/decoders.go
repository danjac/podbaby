package decoders

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/asaskevich/govalidator"
)

type Decoder interface {
	Decode() error
}

var (
	ErrInvalidPassword = errors.New("Password must be at least 6 characters in length")
	ErrInvalidEmail    = errors.New("Invalid email address")
)

// Decode decodes JSON body of request and runs through validator
func Decode(r *http.Request, decoder interface{}) error {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(decoder); err != nil {
		return err
	}
	//return decoder.Decode()
	_, err := govalidator.ValidateStruct(decoder)
	return err
}

type SimpleDecoder struct{}

func (d *SimpleDecoder) Decode() error {
	_, err := govalidator.ValidateStruct(d)
	return err
}

type RecoverPassword struct {
	*SimpleDecoder
	Identifier string `json:"identifier",valid:"required"`
}

type Signup struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Signup) Decode() error {
	if len(s.Name) < 3 {
		return errors.New("Username must be at least 3 characters long")
	}
	if len(s.Password) < 6 {
		return ErrInvalidPassword
	}
	if !govalidator.IsEmail(s.Email) {
		return ErrInvalidEmail
	}
	return nil
}

type Login struct {
	*SimpleDecoder
	Identifier string `json:"identifier" valid:"required"`
	Password   string `json:"password" valid:"required"`
}

type NewChannel struct {
	*SimpleDecoder
	URL string `json:"url" valid:"required,url"`
}

type NewEmail struct {
	*SimpleDecoder
	Email string `json:"email" valid:"required,email"`
}

type NewPassword struct {
	*SimpleDecoder
	OldPassword string `json:"oldPassword"  valid:"required"`
	NewPassword string `json:"newPassword"  valid:"required,length(6|1000)"`
}
