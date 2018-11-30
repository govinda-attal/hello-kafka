package main

import (
	"fmt"

	"github.com/govinda-attal/hello-kafka/pkg/core/status"
)

func main() {

	err := fmt.Errorf("this is error")

	if err != nil {
		errSvc, ok := err.(status.ErrServiceStatus)
		if !ok {
			errSvc = status.ErrInternal.WithError(err)
		}
		err = errSvc
	}
	fmt.Printf("%T", err)
}
