package decoders

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"net/http"
)

// decodes JSON body of request and runs through validator
func decode(r *http.Request, data interface{}) error {
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

func (d *Signup) Decode(r *http.Request) error {
	return decode(r, d)
}

type Login struct {
	Identifier string `json:"identifier",valid:"required"`
	Password   string `json:"password",valid:"required"`
}

func (d *Login) Decode(r *http.Request) error {
	return decode(r, d)
}

type NewChannel struct {
	URL string `json:"url",valid:"url,required"`
}

func (d *NewChannel) Decode(r *http.Request) error {
	return decode(r, d)
}
