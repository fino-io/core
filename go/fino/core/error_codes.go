package core

import "net/http"

var (
	errorCodeIndex     = map[int32]*ErrorCode{}
	errorCodeNameIndex = map[string]*ErrorCode{}
)

func init() {
	errorCodeIndex = map[int32]*ErrorCode{}
	errorCodeNameIndex = map[string]*ErrorCode{}

	addIndex := func(codes ...*ErrorCode) {
		for _, code := range codes {
			errorCodeIndex[code.Code] = code
			errorCodeNameIndex[code.Name] = code
		}
	}

	addIndex(BadRequest, InvalidArgument, MalformedRequest, FailedPrecondition, OutOfRange, Unauthenticated,
		PermissionDenied, NotFound, AlreadyExists, Aborted, ResourceExhausted, Cancelled, UnknownError,
		InternalError, DataLoss, Unimplemented, Unavailable, DeadlineExceeded)
}

var (
	BadRequest = &ErrorCode{
		Code:           400,
		Name:           "BAD_REQUEST",
		Description:    "The request could not be understood by the server due to malformed syntax.",
		HttpStatusCode: http.StatusBadRequest,
	}

	InvalidArgument = &ErrorCode{
		Code:           3,
		Name:           "INVALID_ARGUMENT",
		Description:    "The client specified an invalid argument.",
		HttpStatusCode: http.StatusBadRequest,
	}

	MalformedRequest = &ErrorCode{
		Code:           5,
		Name:           "MALFORMED_SYNTAX",
		Description:    "The syntax of the requested string is malformed.",
		HttpStatusCode: http.StatusBadRequest,
	}

	FailedPrecondition = &ErrorCode{
		Code:           9,
		Name:           "FAILED_PRECONDITION",
		Description:    "The operation was rejected because the system is not in a state required for the operation's execution.",
		HttpStatusCode: http.StatusBadRequest,
	}

	OutOfRange = &ErrorCode{
		Code:           11,
		Name:           "OUT_OF_RANGE",
		Description:    "The operation was attempted past the invalid range.",
		HttpStatusCode: http.StatusBadRequest,
	}

	Unauthenticated = &ErrorCode{
		Code:           401,
		Name:           "UNAUTHENTICATED",
		Description:    "The request does not have the valid authentication credentials for operation.",
		HttpStatusCode: http.StatusUnauthorized,
	}

	PermissionDenied = &ErrorCode{
		Code:           403,
		Name:           "PERMISSION_DENIED",
		Description:    "The caller does not have the permission to execute the specified request.",
		HttpStatusCode: http.StatusForbidden,
	}

	NotFound = &ErrorCode{
		Code:           404,
		Name:           "NOT_FOUND",
		Description:    "Some requested entity (e.g., file or directory) was not found.",
		HttpStatusCode: http.StatusNotFound,
	}

	AlreadyExists = &ErrorCode{
		Code:           6,
		Name:           "ALREADY_EXISTS",
		Description:    "The entity that a client attempted to create (e.g., file or directory) already exists.",
		HttpStatusCode: http.StatusConflict,
	}

	Aborted = &ErrorCode{
		Code:           10,
		Name:           "ABORTED",
		Description:    "The operation was aborted, typically due to a concurrency issue such as a sequencer check failure or transaction abort.",
		HttpStatusCode: http.StatusConflict,
	}

	ResourceExhausted = &ErrorCode{
		Code:           429,
		Name:           "RESOURCE_EXHAUSTED",
		Description:    "Some resource has been exhausted, perhaps a per-user quota, or perhaps the entire file system is out of space.",
		HttpStatusCode: http.StatusTooManyRequests,
	}

	Cancelled = &ErrorCode{
		Code:           499,
		Name:           "CANCELLED",
		Description:    "The operation was cancelled, typically by the caller.",
		HttpStatusCode: 499,
	}

	UnknownError = &ErrorCode{
		Code:           2,
		Name:           "UNKNOWN_ERROR",
		Description:    "Unknown error. For example, this error may be returned when a Status Code received from another address space belongs to an error space that is not known in this address space.",
		HttpStatusCode: http.StatusInternalServerError,
	}

	InternalError = &ErrorCode{
		Code:           500,
		Name:           "INTERNAL_ERROR",
		Description:    "Internal errors. This means that some invariants expected by the underlying system have been broken.",
		HttpStatusCode: http.StatusInternalServerError,
	}

	DataLoss = &ErrorCode{
		Code:           15,
		Name:           "DATA_LOSS",
		Description:    "Unrecoverable data loss or corruption.",
		HttpStatusCode: http.StatusInternalServerError,
	}

	Unimplemented = &ErrorCode{
		Code:           501,
		Name:           "UNIMPLEMENTED",
		Description:    "The operation is not implemented or is not supported/enabled in this service.",
		HttpStatusCode: http.StatusNotImplemented,
	}

	Unavailable = &ErrorCode{
		Code:           503,
		Name:           "UNAVAILABLE",
		Description:    "The service is currently unavailable.",
		HttpStatusCode: http.StatusServiceUnavailable,
	}

	DeadlineExceeded = &ErrorCode{
		Code:           504,
		Name:           "DEADLINE_EXCEEDED",
		Description:    "The deadline expired before the operation could complete.",
		HttpStatusCode: http.StatusGatewayTimeout,
	}
)
