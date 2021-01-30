package errors

func BadRequest(format string, a ...interface{}) error {
	return NewError(400, format, a...)
}

func BadRequestID(id string, format string, a ...interface{}) error {
	return NewErrorID(id, 400, format, a...)
}

func Unauthorized(format string, a ...interface{}) error {
	return NewError(401, format, a...)
}

func UnauthorizedID(id string, format string, a ...interface{}) error {
	return NewErrorID(id, 401, format, a...)
}

func Forbidden(format string, a ...interface{}) error {
	return NewError(403, format, a...)
}

func ForbiddenID(id string, format string, a ...interface{}) error {
	return NewErrorID(id, 403, format, a...)
}

func NotFound(format string, a ...interface{}) error {
	return NewError(404, format, a...)
}

func NotFoundID(id string, format string, a ...interface{}) error {
	return NewErrorID(id, 404, format, a...)
}

func InternalServerError(format string, a ...interface{}) error {
	return NewError(500, format, a...)
}

func InternalServerErrorID(id string, format string, a ...interface{}) error {
	return NewErrorID(id, 500, format, a...)
}
