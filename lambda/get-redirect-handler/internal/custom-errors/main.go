package customerrors

type CustomError struct {
	HTTPStatusCode int
	Message        string
}

func (ce *CustomError) Error() string {
	return ce.Message
}
