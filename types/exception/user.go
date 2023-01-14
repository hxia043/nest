// Package exception define all self exception types.
// This file defines errors which about users.
package exception

type LogoutError struct {
    baseError
}

func NewLogoutError() LogoutError {
    e := LogoutError{}
    e.baseError = newBaseError("logout error")
    return e
}
