package errors

import (
	"fmt"
	"strings"
)

// ==================== juju adaptor start ========================

// Trace just calls AddStack.
func Trace(err error) error {
	if err == nil {
		return nil
	}
	return AddStack(err)
}

// Annotate adds a message and ensures there is a stack trace.
func Annotate(err error, message string) error {
	if err == nil {
		return nil
	}
	hasStack := HasStack(err)
	err = &withMessage{
		cause:         err,
		msg:           message,
		causeHasStack: hasStack,
	}
	if hasStack {
		return err
	}
	return &withStack{
		err,
		callers(),
	}
}

// Annotatef adds a message and ensures there is a stack trace.
func Annotatef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	hasStack := HasStack(err)
	err = &withMessage{
		cause:         err,
		msg:           fmt.Sprintf(format, args...),
		causeHasStack: hasStack,
	}
	if hasStack {
		return err
	}
	return &withStack{
		err,
		callers(),
	}
}

var emptyStack stack

// NewNoStackError creates error without error stack
// later duplicate trace will no longer generate Stack too.
func NewNoStackError(msg string) error {
	return &fundamental{
		msg:   msg,
		stack: &emptyStack,
	}
}

// NewNoStackErrorf creates error with error stack and formats according
// to a format specifier and returns the string as a value that satisfies error.
func NewNoStackErrorf(format string, args ...interface{}) error {
	return &fundamental{
		msg:   fmt.Sprintf(format, args...),
		stack: &emptyStack,
	}
}

// SuspendStack suspends stack generate for error.
// Deprecated, it's semantic is to clear the stack inside, we still allow upper
// layer to add stack again by using Trace.
// Sometimes we have very deep calling stack, the lower layer calls SuspendStack,
// but the upper layer want to add stack to it, if we disable adding stack permanently
// for an error, it's very hard to diagnose certain issues.
func SuspendStack(err error) error {
	if err == nil {
		return err
	}
	cleared := clearStack(err)
	if cleared {
		return err
	}
	return &withStack{
		err,
		&emptyStack,
	}
}

func clearStack(err error) (cleared bool) {
	switch typedErr := err.(type) {
	case *withMessage:
		return clearStack(typedErr.Cause())
	case *fundamental:
		typedErr.stack = &emptyStack
		return true
	case *withStack:
		typedErr.stack = &emptyStack
		clearStack(typedErr.Cause())
		return true
	default:
		return false
	}
}

// ErrorStack will format a stack trace if it is available, otherwise it will be Error()
// If the error is nil, the empty string is returned
// Note that this just calls fmt.Sprintf("%+v", err)
func ErrorStack(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("%+v", err)
}

// IsNotFound reports whether err was not found error.
func IsNotFound(err error) bool {
	return strings.Contains(err.Error(), "not found")
}

// NotFoundf represents an error with not found message.
func NotFoundf(format string, args ...interface{}) error {
	return Errorf(format+" not found", args...)
}

// BadRequestf represents an error with bad request message.
func BadRequestf(format string, args ...interface{}) error {
	return Errorf(format+" bad request", args...)
}

// NotSupportedf represents an error with not supported message.
func NotSupportedf(format string, args ...interface{}) error {
	return Errorf(format+" not supported", args...)
}

// NotValidf represents an error with not valid message.
func NotValidf(format string, args ...interface{}) error {
	return Errorf(format+" not valid", args...)
}

// IsAlreadyExists reports whether err was already exists error.
func IsAlreadyExists(err error) bool {
	return strings.Contains(err.Error(), "already exists")
}

// AlreadyExistsf represents an error with already exists message.
func AlreadyExistsf(format string, args ...interface{}) error {
	return Errorf(format+" already exists", args...)
}

// ==================== juju adaptor end ========================
