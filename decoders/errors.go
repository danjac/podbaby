package decoders

import (
	"fmt"
	"reflect"

	"github.com/asaskevich/govalidator"
)

// store the errors as a map
type Errors map[string]string

func (e Errors) Error() string {
	return fmt.Sprintf("%v", e)
}

func makeValidationErrors(src govalidator.Errors, data interface{}) Errors {
	// use JSON tags as error keys
	if src == nil {
		return nil
	}

	dst := make(Errors)

	value := reflect.ValueOf(data).Elem()
	t := value.Type()

	for i := 0; i < value.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("json")
		if tag == "" {
			tag = f.Name
		}
		for _, err := range src {
			vErr, ok := err.(govalidator.Error)
			if ok {
				if vErr.Name == f.Name {
					dst[tag] = vErr.Err.Error()
				}
			}
		}
	}
	return dst
}
