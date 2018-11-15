package valderr

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/govinda-attal/hello-kafka/pkg/core/status"
)

func NewErrStatusWithValErrors(e status.ErrServiceStatus, valErrs validation.Errors) status.ErrServiceStatus {
	errSvc := status.ErrServiceStatus{ServiceStatus: status.ServiceStatus{Code: e.Code, Message: e.Message, Details: nil}}
	for _, msg := range valErrs {
		errSvc.AddDtlMsg(msg.Error())
	}
	return errSvc
}
