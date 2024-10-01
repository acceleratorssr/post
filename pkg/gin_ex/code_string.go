package gin_ex

import "strconv"

func (c Code) String() string {
	switch c {
	case OK:
		return "OK"
	case Canceled:
		return "Canceled"
	case Unknown:
		return "Unknown"
	case InvalidArgument:
		return "InvalidArgument"
	case DeadlineExceeded:
		return "DeadlineExceeded"
	case NotFound:
		return "NotFound"
	case AlreadyExists:
		return "AlreadyExists"
	case PermissionDenied:
		return "PermissionDenied"
	case ResourceExhausted:
		return "ResourceExhausted"
	case FailedPrecondition:
		return "FailedPrecondition"
	case Aborted:
		return "Aborted"
	case Internal:
		return "Internal"
	case Unavailable:
		return "Unavailable"
	case Unauthenticated:
		return "Unauthenticated"
	case System:
		return "System"
	default:
		return "Code(" + strconv.FormatInt(int64(c), 10) + ")"
	}
}

func canonicalString(c Code) string {
	switch c {
	case OK:
		return "OK"
	case Canceled:
		return "CANCELLED"
	case Unknown:
		return "UNKNOWN"
	case InvalidArgument:
		return "INVALID_ARGUMENT"
	case DeadlineExceeded:
		return "DEADLINE_EXCEEDED"
	case NotFound:
		return "NOT_FOUND"
	case AlreadyExists:
		return "ALREADY_EXISTS"
	case PermissionDenied:
		return "PERMISSION_DENIED"
	case ResourceExhausted:
		return "RESOURCE_EXHAUSTED"
	case FailedPrecondition:
		return "FAILED_PRECONDITION"
	case Aborted:
		return "ABORTED"
	case Internal:
		return "INTERNAL"
	case Unavailable:
		return "UNAVAILABLE"
	case Unauthenticated:
		return "UNAUTHENTICATED"
	case System:
		return "SYSTEM"
	default:
		return "CODE(" + strconv.FormatInt(int64(c), 10) + ")"
	}
}
