package errors

import "fmt"

type HttpStatusCode uint32

type VnfdsvcError struct {
	OrigError error
	HttpCode  HttpStatusCode
}

func (e *VnfdsvcError) String() string {
	return fmt.Sprintf("%s, HttpStatusCode=%d", e.OrigError.Error(), e.HttpCode)
}

func (e *VnfdsvcError) Error() string {
	return e.String()
}
