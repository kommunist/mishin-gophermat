package notfound

type NotFoundError struct {
	Err error
}

func (e *NotFoundError) Error() string {
	return e.Err.Error()
}

func NewNotFoundError(err error) error {
	return &NotFoundError{Err: err}
}
