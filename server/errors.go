package server

type Error interface {
	error
	Status() int
}

type HTTPError struct {
	Code int
	Err  error
}

func (e HTTPError) Error() string {
	return e.Err.Error()
}

func (e HTTPError) Status() int {
	return e.Code
}
