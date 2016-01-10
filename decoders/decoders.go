package decoders

import (
	"encoding/json"
	"net/http"

	"github.com/asaskevich/govalidator"
)

type Decoder interface {
	Decode() error
}

// Decode decodes JSON body of request and runs through validator
func Decode(r *http.Request, data interface{}) error {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return err
	}
	// decoder does some more in depth validation
	if decoder, ok := data.(Decoder); ok {
		return decoder.Decode()
	}
	// just do some basic validation
	if _, err := govalidator.ValidateStruct(data); err != nil {
		if vErrors, ok := err.(govalidator.Errors); ok {
			return makeValidationErrors(vErrors, data)
		}
		return err
	}
	return nil
}

type RecoverPassword struct {
	Identifier string `json:"identifier",valid:"required"`
}

type Signup struct {
	Name     string `json:"name" valid:"required,length(3|60)"`
	Email    string `json:"email" valid:"required,email"`
	Password string `json:"password" valid:"required,length(6|1000)"`
}

type Login struct {
	Identifier string `json:"identifier" valid:"required"`
	Password   string `json:"password" valid:"required"`
}

type NewChannel struct {
	URL string `json:"url" valid:"required,url"`
}

type NewEmail struct {
	Email string `json:"email" valid:"required,email"`
}

type NewPassword struct {
	OldPassword string `json:"oldPassword"  valid:"required"`
	NewPassword string `json:"newPassword"  valid:"required,length(6|1000)"`
}
