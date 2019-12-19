package mongo

import (
	"fmt"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

func newError(cause error, code int32, message string) *Error {
	return &Error{
		Code:     code,
		CodeName: CodeToName(code),
		Message:  message,
		Cause:    cause,
	}
}

func newErrorf(cause error, code int32, format string, args ...interface{}) *Error {
	return newError(cause, code, fmt.Sprintf(format, args...))
}

// Error is a special error used to return an error message to a client.
type Error struct {
	Code     int32
	CodeName string
	Message  string

	Cause error
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}

	return e.Message
}

// MarshalBSON implements the bson.Marshaler interface.
func (e *Error) MarshalBSON() []byte {
	_, doc := bsoncore.AppendDocumentStart(nil)
	doc = bsoncore.AppendInt32Element(doc, "ok", 0)
	doc = bsoncore.AppendStringElement(doc, "errmsg", e.Message)

	if code := e.Code; code != 0 {
		doc = bsoncore.AppendInt32Element(doc, "code", code)
	}

	if codeName := e.CodeName; codeName != "" {
		doc = bsoncore.AppendStringElement(doc, "codeName", codeName)
	}

	doc, _ = bsoncore.AppendDocumentEnd(doc, 0)

	return doc
}

// Unwrap is used by errors.Unwrap.
func (e *Error) Unwrap() error {
	return e.Cause
}
