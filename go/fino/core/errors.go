package core

import (
	"errors"
)

type basicError Error

func newBasicError(code *ErrorCode, message string, arguments ...any) *basicError {
	code = NewErrorCode(code.GetCode())
	if len(arguments) == 0 {
		return (*basicError)(NewError(code, message))
	}
	return (*basicError)(NewErrorf(code, message, arguments...))
}

func (e *basicError) Error() string {
	return (*Error)(e).Error()
}

func (e *basicError) ToError() *Error {
	return (*Error)(e)
}

// Unwrap exposes the shared protobuf error to errors.Is and errors.As.
func (e *basicError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.ToError()
}

func (e *basicError) StatusCode() int {
	return (*Error)(e).StatusCode()
}

func (e *basicError) AddDetail(detail any) *basicError {
	return (*basicError)((*Error)(e).AddDetail(detail))
}

type BadRequestError struct {
	*basicError
}

func NewBadRequestError(format string, args ...any) *BadRequestError {
	return &BadRequestError{newBasicError(BadRequest, format, args...)}
}

func IsBadRequestError(err error) bool {
	return errors.Is(err, &BadRequestError{})
}

func (*BadRequestError) Is(target error) bool {
	_, ok := target.(*BadRequestError)
	return ok
}

type InvalidArgumentError struct {
	*basicError
}

func NewInvalidArgumentError(format string, args ...any) *InvalidArgumentError {
	return &InvalidArgumentError{newBasicError(InvalidArgument, format, args...)}
}

func IsInvalidArgumentError(err error) bool {
	return errors.Is(err, &InvalidArgumentError{})
}

func (*InvalidArgumentError) Is(target error) bool {
	_, ok := target.(*InvalidArgumentError)
	return ok
}

type MalformedRequestError struct {
	*basicError
}

func NewMalformedRequestError(format string, args ...any) *MalformedRequestError {
	return &MalformedRequestError{newBasicError(MalformedRequest, format, args...)}
}

func IsMalformedRequestError(err error) bool {
	return errors.Is(err, &MalformedRequestError{})
}

func (*MalformedRequestError) Is(target error) bool {
	_, ok := target.(*MalformedRequestError)
	return ok
}

type FailedPreconditionError struct {
	*basicError
}

func NewFailedPreconditionError(format string, args ...any) *FailedPreconditionError {
	return &FailedPreconditionError{newBasicError(FailedPrecondition, format, args...)}
}

func IsFailedPreconditionError(err error) bool {
	return errors.Is(err, &FailedPreconditionError{})
}

func (*FailedPreconditionError) Is(target error) bool {
	_, ok := target.(*FailedPreconditionError)
	return ok
}

type OutOfRangeError struct {
	*basicError
}

func NewOutOfRangeError(format string, args ...any) *OutOfRangeError {
	return &OutOfRangeError{newBasicError(OutOfRange, format, args...)}
}

func IsOutOfRangeError(err error) bool {
	return errors.Is(err, &OutOfRangeError{})
}

func (*OutOfRangeError) Is(target error) bool {
	_, ok := target.(*OutOfRangeError)
	return ok
}

type UnauthenticatedError struct {
	*basicError
}

func NewUnauthenticatedError(format string, args ...any) *UnauthenticatedError {
	return &UnauthenticatedError{newBasicError(Unauthenticated, format, args...)}
}

func IsUnauthenticatedError(err error) bool {
	return errors.Is(err, &UnauthenticatedError{})
}

func (*UnauthenticatedError) Is(target error) bool {
	_, ok := target.(*UnauthenticatedError)
	return ok
}

type PermissionDeniedError struct {
	*basicError
}

func NewPermissionDeniedError(format string, args ...any) *PermissionDeniedError {
	return &PermissionDeniedError{newBasicError(PermissionDenied, format, args...)}
}

func IsPermissionDeniedError(err error) bool {
	return errors.Is(err, &PermissionDeniedError{})
}

func (*PermissionDeniedError) Is(target error) bool {
	_, ok := target.(*PermissionDeniedError)
	return ok
}

type NotFoundError struct {
	*basicError
}

func NewNotFoundError(format string, args ...any) *NotFoundError {
	return &NotFoundError{newBasicError(NotFound, format, args...)}
}

func IsNotFoundError(err error) bool {
	return errors.Is(err, &NotFoundError{})
}

func (*NotFoundError) Is(target error) bool {
	_, ok := target.(*NotFoundError)
	return ok
}

type AlreadyExistsError struct {
	*basicError
}

func NewAlreadyExistsError(format string, args ...any) *AlreadyExistsError {
	return &AlreadyExistsError{newBasicError(AlreadyExists, format, args...)}
}

func IsAlreadyExistsError(err error) bool {
	return errors.Is(err, &AlreadyExistsError{})
}

func (*AlreadyExistsError) Is(target error) bool {
	_, ok := target.(*AlreadyExistsError)
	return ok
}

type AbortedError struct {
	*basicError
}

func NewAbortedError(format string, args ...any) *AbortedError {
	return &AbortedError{newBasicError(Aborted, format, args...)}
}

func IsAbortedError(err error) bool {
	return errors.Is(err, &AbortedError{})
}

func (*AbortedError) Is(target error) bool {
	_, ok := target.(*AbortedError)
	return ok
}

type ResourceExhaustedError struct {
	*basicError
}

func NewResourceExhaustedError(format string, args ...any) *ResourceExhaustedError {
	return &ResourceExhaustedError{newBasicError(ResourceExhausted, format, args...)}
}

func IsResourceExhaustedError(err error) bool {
	return errors.Is(err, &ResourceExhaustedError{})
}

func (*ResourceExhaustedError) Is(target error) bool {
	_, ok := target.(*ResourceExhaustedError)
	return ok
}

type CancelledError struct {
	*basicError
}

func NewCancelledError(format string, args ...any) *CancelledError {
	return &CancelledError{newBasicError(Cancelled, format, args...)}
}

func IsCancelledError(err error) bool {
	return errors.Is(err, &CancelledError{})
}

func (*CancelledError) Is(target error) bool {
	_, ok := target.(*CancelledError)
	return ok
}

type UnknownErrorError struct {
	*basicError
}

func NewUnknownErrorError(format string, args ...any) *UnknownErrorError {
	return &UnknownErrorError{newBasicError(UnknownError, format, args...)}
}

func IsUnknownErrorError(err error) bool {
	return errors.Is(err, &UnknownErrorError{})
}

func (*UnknownErrorError) Is(target error) bool {
	_, ok := target.(*UnknownErrorError)
	return ok
}

type InternalErrorError struct {
	*basicError
}

func NewInternalErrorError(format string, args ...any) *InternalErrorError {
	return &InternalErrorError{newBasicError(InternalError, format, args...)}
}

func IsInternalError(err error) bool {
	return errors.Is(err, &InternalErrorError{})
}

func (*InternalErrorError) Is(target error) bool {
	_, ok := target.(*InternalErrorError)
	return ok
}

type DataLossError struct {
	*basicError
}

func NewDataLossError(format string, args ...any) *DataLossError {
	return &DataLossError{newBasicError(DataLoss, format, args...)}
}

func IsDataLossError(err error) bool {
	return errors.Is(err, &DataLossError{})
}

func (*DataLossError) Is(target error) bool {
	_, ok := target.(*DataLossError)
	return ok
}

type UnimplementedError struct {
	*basicError
}

func NewUnimplementedError(format string, args ...any) *UnimplementedError {
	return &UnimplementedError{newBasicError(Unimplemented, format, args...)}
}

func IsUnimplementedError(err error) bool {
	return errors.Is(err, &UnimplementedError{})
}

func (*UnimplementedError) Is(target error) bool {
	_, ok := target.(*UnimplementedError)
	return ok
}

type UnavailableError struct {
	*basicError
}

func NewUnavailableError(format string, args ...any) *UnavailableError {
	return &UnavailableError{newBasicError(Unavailable, format, args...)}
}

func IsUnavailableError(err error) bool {
	return errors.Is(err, &UnavailableError{})
}

func (*UnavailableError) Is(target error) bool {
	_, ok := target.(*UnavailableError)
	return ok
}

type DeadlineExceededError struct {
	*basicError
}

func NewDeadlineExceededError(format string, args ...any) *DeadlineExceededError {
	return &DeadlineExceededError{newBasicError(DeadlineExceeded, format, args...)}
}

func IsDeadlineExceededError(err error) bool {
	return errors.Is(err, &DeadlineExceededError{})
}

func (*DeadlineExceededError) Is(target error) bool {
	_, ok := target.(*DeadlineExceededError)
	return ok
}
