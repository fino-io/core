package core

import (
	"errors"
)

type basicError Error

func newBasicError(code *ErrorCode, message string, arguments ...any) *basicError {
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

func (*BadRequestError) Is(err error) bool {
	var badRequestError *BadRequestError
	return errors.As(err, &badRequestError)
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

func (*InvalidArgumentError) Is(err error) bool {
	var badRequestError *InvalidArgumentError
	return errors.As(err, &badRequestError)
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

func (*MalformedRequestError) Is(err error) bool {
	var malformedRequestError *MalformedRequestError
	return errors.As(err, &malformedRequestError)
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

func (*FailedPreconditionError) Is(err error) bool {
	var failedPreconditionError *FailedPreconditionError
	return errors.As(err, &failedPreconditionError)
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

func (*OutOfRangeError) Is(err error) bool {
	var outOfRangeError *OutOfRangeError
	return errors.As(err, &outOfRangeError)
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

func (*UnauthenticatedError) Is(err error) bool {
	var unauthenticatedError *UnauthenticatedError
	return errors.As(err, &unauthenticatedError)
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

func (*PermissionDeniedError) Is(err error) bool {
	var permissionDeniedError *PermissionDeniedError
	return errors.As(err, &permissionDeniedError)
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

func (*NotFoundError) Is(err error) bool {
	var notFoundError *NotFoundError
	return errors.As(err, &notFoundError)

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

func (*AlreadyExistsError) Is(err error) bool {
	var alreadyExistsError *AlreadyExistsError
	return errors.As(err, &alreadyExistsError)
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

func (*AbortedError) Is(err error) bool {
	var abortedError *AbortedError
	return errors.As(err, &abortedError)
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

func (*ResourceExhaustedError) Is(err error) bool {
	var resourceExhaustedError *ResourceExhaustedError
	return errors.As(err, &resourceExhaustedError)
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

func (*CancelledError) Is(err error) bool {
	var cancelledError *CancelledError
	return errors.As(err, &cancelledError)
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

func (*UnknownErrorError) Is(err error) bool {
	var unknownErrorError *UnknownErrorError
	return errors.As(err, &unknownErrorError)
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

func (*InternalErrorError) Is(err error) bool {
	var internalErrorErrorError *InternalErrorError
	return errors.As(err, &internalErrorErrorError)
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

func (*DataLossError) Is(err error) bool {
	var dataLossError *DataLossError
	return errors.As(err, &dataLossError)
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

func (*UnimplementedError) Is(err error) bool {
	var unimplementedError *UnimplementedError
	return errors.As(err, &unimplementedError)
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

func (*UnavailableError) Is(err error) bool {
	var unavailableErrorError *UnavailableError
	return errors.As(err, &unavailableErrorError)
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

func (*DeadlineExceededError) Is(err error) bool {
	var deadlineExceededError *DeadlineExceededError
	return errors.As(err, &deadlineExceededError)
}
