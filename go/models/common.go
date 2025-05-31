package models

import (
	"time"

	"gorm.io/gorm"
)

// StatusCode custom type for representing status codes
type StatusCode string

// Constants for status codes
const (
	StatusOK                   StatusCode = "OK"
	StatusCreated              StatusCode = "CREATED"
	StatusAccepted             StatusCode = "ACCEPTED"
	StatusNonAuthoritativeInfo StatusCode = "NON_AUTHORITATIVE_INFORMATION"
	StatusNoContent            StatusCode = "NO_CONTENT"
	StatusResetContent         StatusCode = "RESET_CONTENT"
	StatusPartialContent       StatusCode = "PARTIAL_CONTENT"
	StatusMovedPermanently     StatusCode = "MOVED_PERMANENTLY"
	StatusFound                StatusCode = "FOUND"
	StatusSeeOther             StatusCode = "SEE_OTHER"
	StatusNotModified          StatusCode = "NOT_MODIFIED"
	StatusTemporaryRedirect    StatusCode = "TEMPORARY_REDIRECT"
	StatusPermanentRedirect    StatusCode = "PERMANENT_REDIRECT"
	StatusBadRequest           StatusCode = "BAD_REQUEST"
	StatusUnauthorized         StatusCode = "UNAUTHORIZED"
	StatusPaymentRequired      StatusCode = "PAYMENT_REQUIRED"
	StatusForbidden            StatusCode = "FORBIDDEN"
	StatusNotFound             StatusCode = "NOT_FOUND"
	StatusMethodNotAllowed     StatusCode = "METHOD_NOT_ALLOWED"
	StatusNotAcceptable        StatusCode = "NOT_ACCEPTABLE"
	StatusProxyAuthRequired    StatusCode = "PROXY_AUTHENTICATION_REQUIRED"
	StatusRequestTimeout       StatusCode = "REQUEST_TIMEOUT"
	StatusConflict             StatusCode = "CONFLICT"
	StatusGone                 StatusCode = "GONE"
	StatusLengthRequired       StatusCode = "LENGTH_REQUIRED"
	StatusPreconditionFailed   StatusCode = "PRECONDITION_FAILED"
	StatusRequestEntityTooLarge StatusCode = "REQUEST_ENTITY_TOO_LARGE"
	StatusRequestURITooLong    StatusCode = "REQUEST_URI_TOO_LONG"
	StatusUnsupportedMediaType StatusCode = "UNSUPPORTED_MEDIA_TYPE"
	StatusRequestedRangeNotSatisfiable StatusCode = "REQUESTED_RANGE_NOT_SATISFIABLE"
	StatusExpectationFailed    StatusCode = "EXPECTATION_FAILED"
	StatusTeapot               StatusCode = "I_AM_A_TEAPOT"
	StatusUnprocessableEntity  StatusCode = "UNPROCESSABLE_ENTITY"
	StatusLocked               StatusCode = "LOCKED"
	StatusFailedDependency     StatusCode = "FAILED_DEPENDENCY"
	StatusTooEarly             StatusCode = "TOO_EARLY"
	StatusUpgradeRequired      StatusCode = "UPGRADE_REQUIRED"
	StatusPreconditionRequired StatusCode = "PRECONDITION_REQUIRED"
	StatusTooManyRequests      StatusCode = "TOO_MANY_REQUESTS"
	StatusRequestHeaderFieldsTooLarge StatusCode = "REQUEST_HEADER_FIELDS_TOO_LARGE"
	StatusUnavailableForLegalReasons StatusCode = "UNAVAILABLE_FOR_LEGAL_REASONS"
	StatusInternalServerError  StatusCode = "INTERNAL_SERVER_ERROR"
	StatusNotImplemented       StatusCode = "NOT_IMPLEMENTED"
	StatusBadGateway           StatusCode = "BAD_GATEWAY"
	StatusServiceUnavailable   StatusCode = "SERVICE_UNAVAILABLE"
	StatusGatewayTimeout       StatusCode = "GATEWAY_TIMEOUT"
	StatusHTTPVersionNotSupported StatusCode = "HTTP_VERSION_NOT_SUPPORTED"
	StatusVariantAlsoNegotiates StatusCode = "VARIANT_ALSO_NEGOTIATES"
	StatusInsufficientStorage  StatusCode = "INSUFFICIENT_STORAGE"
	StatusLoopDetected         StatusCode = "LOOP_DETECTED"
	StatusNotExtended          StatusCode = "NOT_EXTENDED"
	StatusNetworkAuthenticationRequired StatusCode = "NETWORK_AUTHENTICATION_REQUIRED"
)

// BaseModel defines common fields for GORM models
type BaseModel struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
