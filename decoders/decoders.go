package decoders

import (
	"encoding/json"
	"net/http"

	"github.com/asaskevich/govalidator"
)

// decodes JSON body of request and runs through validator
func Decode(r *http.Request, data interface{}) error {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return err
	}
	if _, err := govalidator.ValidateStruct(data); err != nil {
		return err
	}
	return nil
}

type Signup struct {
	Name     string `json:"name",valid:"required"`
	Email    string `json:"email",valid:"email,required"`
	Password string `json:"password",valid:"required"`
}

type Login struct {
	Identifier string `json:"identifier",valid:"required"`
	Password   string `json:"password",valid:"required"`
}

type NewChannel struct {
	URL string `json:"url",valid:"url,required"`
}

type NewEmail struct {
	Email string `json:"email",valid:"email,required"`
}
