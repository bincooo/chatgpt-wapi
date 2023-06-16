package chatgpt

import "strconv"

type Error struct {
	Code    int
	Message string
}

func (c Error) Error() string {
	return "error { code=" + strconv.Itoa(c.Code) + ", message=" + c.Message + " }"
}

func NewError(code int, message string) error {
	return Error{code, message}
}
