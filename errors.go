package main

import (
	"fmt"
)

type ErrUnsolvable struct {
	msg string
}

func NewErrUnsolvable(format string, a ...interface{}) ErrUnsolvable {
	return ErrUnsolvable{msg: fmt.Sprintf(format, a...)}
}

func (e ErrUnsolvable) Error() string {
	return e.msg
}
