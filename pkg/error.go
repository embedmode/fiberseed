package pkg

type Error struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func EntityNotFound(m string) *Error {
	return &Error{Status: 404, Code: "entity-not-found", Message: m}
}

func BadRequest(m string) *Error {
	return &Error{Status: 400, Code: "bad-request", Message: m}
}
