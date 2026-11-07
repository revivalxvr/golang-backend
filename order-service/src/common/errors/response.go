package errors

type Response struct {
	Message  string
	HttpCode int
}

func (r Response) Error() string {
	return r.Message
}
