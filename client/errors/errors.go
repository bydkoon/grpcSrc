package errors

func New(text string, step string, err error) *Errors {
	return &Errors{text, step, err}
}

type Errors struct {
	s     string
	step  string
	error error
}

func (e *Errors) ErrorString() string {
	return e.s
}

func (e *Errors) ErrorStep() string {
	return e.step
}

func (e *Errors) Error() error {
	return e.error
}
