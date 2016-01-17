package server

type HTTPError interface {
	error
	Status() int
}

type httpError struct {
	err  error
	code int
}

func (e httpError) Error() string {
	return e.err.Error()
}

func (e httpError) Status() int {
	return e.code
}
