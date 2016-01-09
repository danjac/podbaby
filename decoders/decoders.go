package decoders

import (
	"encoding/json"
	"net/http"

	"github.com/asaskevich/govalidator"
)

// Decode decodes JSON body of request and runs through validator
func Decode(r *http.Request, decoder interface{}) error {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(decoder); err != nil {
		return err
	}
	// tbd: check type for a decoder interface
	_, err := govalidator.ValidateStruct(decoder)
	return err
}

type RecoverPassword struct {
	Identifier string `json:"identifier",valid:"required"`
}

type Signup struct {
	Name     string `json:"name" valid:"required,length(3)"`
	Email    string `json:"email" valid:"required,email"`
	Password string `json:"password" valid:"required,length(6)"`
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
	NewPassword string `json:"newPassword"  valid:"required,length(6)"`
}
