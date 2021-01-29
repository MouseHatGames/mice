package errors

func BadRequest(format string, a ...interface{}) error {
	return NewError(400, format, a...)
}

func Unauthorized(format string, a ...interface{}) error {
	return NewError(401, format, a...)
}

func Forbidden(format string, a ...interface{}) error {
	return NewError(403, format, a...)
}

func NotFound(format string, a ...interface{}) error {
	return NewError(404, format, a...)
}

func InternalServerError(format string, a ...interface{}) error {
	return NewError(404, format, a...)
}
