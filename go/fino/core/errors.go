package core

import "errors"

func NewBadRequestError(format string, args ...any) *Error {
	return NewErrorf(BadRequest, format, args...)
}

func NewInvalidArgumentError(format string, args ...any) *Error {
	return NewErrorf(InvalidArgument, format, args...)
}

func NewMalformedRequestError(format string, args ...any) *Error {
	return NewErrorf(MalformedRequest, format, args...)
}

func NewFailedPreconditionError(format string, args ...any) *Error {
	return NewErrorf(FailedPrecondition, format, args...)
}

func NewOutOfRangeError(format string, args ...any) *Error {
	return NewErrorf(OutOfRange, format, args...)
}

func NewUnauthenticatedError(format string, args ...any) *Error {
	return NewErrorf(Unauthenticated, format, args...)
}

func NewPermissionDeniedError(format string, args ...any) *Error {
	return NewErrorf(PermissionDenied, format, args...)
}

func NewNotFoundError(format string, args ...any) *Error {
	return NewErrorf(NotFound, format, args...)
}

func NewAlreadyExistsError(format string, args ...any) *Error {
	return NewErrorf(AlreadyExists, format, args...)
}

func NewAbortedError(format string, args ...any) *Error {
	return NewErrorf(Aborted, format, args...)
}

func NewResourceExhaustedError(format string, args ...any) *Error {
	return NewErrorf(ResourceExhausted, format, args...)
}

func NewCancelledError(format string, args ...any) *Error {
	return NewErrorf(Cancelled, format, args...)
}

func NewUnknownError(format string, args ...any) *Error {
	return NewErrorf(UnknownError, format, args...)
}

func NewInternalError(format string, args ...any) *Error {
	return NewErrorf(InternalError, format, args...)
}

func NewDataLossError(format string, args ...any) *Error {
	return NewErrorf(DataLoss, format, args...)
}

func NewUnimplementedError(format string, args ...any) *Error {
	return NewErrorf(Unimplemented, format, args...)
}

func NewUnavailableError(format string, args ...any) *Error {
	return NewErrorf(Unavailable, format, args...)
}

func NewDeadlineExceededError(format string, args ...any) *Error {
	return NewErrorf(DeadlineExceeded, format, args...)
}

func IsBadRequestError(err error) bool {
	return errors.Is(err, &Error{Code: BadRequest})
}

func IsInvalidArgumentError(err error) bool {
	return errors.Is(err, &Error{Code: InvalidArgument})
}

func IsMalformedRequestError(err error) bool {
	return errors.Is(err, &Error{Code: MalformedRequest})
}

func IsFailedPreconditionError(err error) bool {
	return errors.Is(err, &Error{Code: FailedPrecondition})
}

func IsOutOfRangeError(err error) bool {
	return errors.Is(err, &Error{Code: OutOfRange})
}

func IsUnauthenticatedError(err error) bool {
	return errors.Is(err, &Error{Code: Unauthenticated})
}

func IsPermissionDeniedError(err error) bool {
	return errors.Is(err, &Error{Code: PermissionDenied})
}

func IsNotFoundError(err error) bool {
	return errors.Is(err, &Error{Code: NotFound})
}

func IsAlreadyExistsError(err error) bool {
	return errors.Is(err, &Error{Code: AlreadyExists})
}

func IsAbortedError(err error) bool {
	return errors.Is(err, &Error{Code: Aborted})
}

func IsResourceExhaustedError(err error) bool {
	return errors.Is(err, &Error{Code: ResourceExhausted})
}

func IsCancelledError(err error) bool {
	return errors.Is(err, &Error{Code: Cancelled})
}

func IsUnknownError(err error) bool {
	return errors.Is(err, &Error{Code: UnknownError})
}

func IsInternalError(err error) bool {
	return errors.Is(err, &Error{Code: InternalError})
}

func IsDataLossError(err error) bool {
	return errors.Is(err, &Error{Code: DataLoss})
}

func IsUnimplementedError(err error) bool {
	return errors.Is(err, &Error{Code: Unimplemented})
}

func IsUnavailableError(err error) bool {
	return errors.Is(err, &Error{Code: Unavailable})
}

func IsDeadlineExceededError(err error) bool {
	return errors.Is(err, &Error{Code: DeadlineExceeded})
}
