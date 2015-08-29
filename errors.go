package asapp

import "errors"

type ErrorDetails map[string]string

type CompoundError interface {
	Error() string
	Details() ErrorDetails
}

type compoundError struct {
	error
	details ErrorDetails
}

func (c compoundError) Details() ErrorDetails {
	return c.details
}

type UserError struct{ compoundError }

type ServerError struct{ compoundError }

type FormError struct{ compoundError }

type AuthenticationError struct{ compoundError }

func NewUserError(msg string) UserError {
	return UserError{compoundError{errors.New(msg), nil}}
}

func NewServerError(msg string) ServerError {
	return ServerError{compoundError{errors.New(msg), nil}}
}

func NewFormError(msg string, details ErrorDetails) FormError {
	return FormError{
		compoundError{
			errors.New(msg),
			details,
		},
	}
}

func NewAuthenticationError(msg string) AuthenticationError {
	return AuthenticationError{compoundError{errors.New(msg), nil}}
}
