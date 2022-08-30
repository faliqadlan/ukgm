package help

import (
	"fmt"

	"github.com/go-errors/errors"
)

// Stacktrace return void and track the error
func Stacktrace(err error) {
	err = errors.Errorf(err.Error())
	fmt.Println(err.(*errors.Error).ErrorStack())
}
