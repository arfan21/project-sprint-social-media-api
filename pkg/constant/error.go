package constant

import (
	"errors"
	"net/http"
)

const (
	ErrSQLUniqueViolation = "23505"
	ErrSQLInvalidUUID     = "22P02"
)

var (
	ErrEmailOrPhoneAlreadyRegistered = &ErrWithCode{HTTPStatusCode: http.StatusConflict, Message: "email or phone already registered"}
	ErrUsernameOrPasswordInvalid     = &ErrWithCode{HTTPStatusCode: http.StatusBadRequest, Message: "username or password invalid"}
	ErrUserNotFound                  = &ErrWithCode{HTTPStatusCode: http.StatusNotFound, Message: "user not found"}
	ErrInvalidUUID                   = errors.New("invalid uuid length or format")
	ErrAccessForbidden               = &ErrWithCode{HTTPStatusCode: http.StatusForbidden, Message: "access forbidden"}
	ErrFriendSelfAdding              = &ErrWithCode{HTTPStatusCode: http.StatusBadRequest, Message: "cannot add self as friend"}
	ErrFriendUseralreadyAdded        = &ErrWithCode{HTTPStatusCode: http.StatusBadRequest, Message: "user already added as friend"}
	ErrFriendSelfDeleting            = &ErrWithCode{HTTPStatusCode: http.StatusBadRequest, Message: "cannot delete self as friend"}
	ErrFriendUserNotAdded            = &ErrWithCode{HTTPStatusCode: http.StatusBadRequest, Message: "user not found in friend list"}
	ErrPostNotFound                  = &ErrWithCode{HTTPStatusCode: http.StatusNotFound, Message: "post not found"}
	ErrUserNotFriend                 = &ErrWithCode{HTTPStatusCode: http.StatusBadRequest, Message: "user is not friend with post owner"}
)

type ErrWithCode struct {
	HTTPStatusCode int
	Message        string
}

func (e *ErrWithCode) Error() string {
	return e.Message
}

type ErrValidation struct {
	Message string
}

func (e *ErrValidation) Error() string {
	return e.Message
}
