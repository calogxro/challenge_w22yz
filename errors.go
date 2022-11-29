package main

const MSG_KEY_EXISTS = "Key exists"
const MSG_KEY_NOTFOUND = "Key not found"

type KeyExists struct{}

func (m *KeyExists) Error() string {
	return MSG_KEY_EXISTS
}

type KeyNotFound struct{}

func (m *KeyNotFound) Error() string {
	return MSG_KEY_NOTFOUND
}
