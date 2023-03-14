package domain

import "errors"

const MSG_KEY_EXISTS = "key exists"
const MSG_KEY_NOTFOUND = "key not found"
const MSG_INPUT_NOTVALID = "input not valid"

// ErrKeyExists is returned when a record already exists for a given key.
var ErrKeyExists = errors.New(MSG_KEY_EXISTS)

// ErrKeyNotFound is returned when a requested record is not found.
var ErrKeyNotFound = errors.New(MSG_KEY_NOTFOUND)

// ErrInputNotValid is returned when not valid input has been submitted.
var ErrInputNotValid = errors.New(MSG_INPUT_NOTVALID)

// type KeyExists struct{}

// func (m *KeyExists) Error() string {
// 	return MSG_KEY_EXISTS
// }

// type KeyNotFound struct{}

// func (m *KeyNotFound) Error() string {
// 	return MSG_KEY_NOTFOUND
// }
