package apierrors

import "errors"

var ErrAuthUserNotFound = errors.New("user not found. Unauthorized")
