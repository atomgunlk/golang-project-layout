package be_error

import (
	"fmt"

	"github.com/pkg/errors"
)

type BeError struct {
	code string
	Err  error
}

func (r *BeError) Code() string {
	return r.code
}
func (r *BeError) Error() string {
	return fmt.Sprintf("%v", r.Err)
}
func (r *BeError) Message() string {
	return fmt.Sprintf("%v", r.Err)
}
func (r *BeError) WrapMessage(m string) string {
	return fmt.Sprintf("%v, %s", r.Err, m)
}
func (r *BeError) FullError() string {
	return fmt.Sprintf("code %s: err %v", r.code, r.Err)
}

var (
	Success = BeError{"00000", errors.New("Success")}

	// Common 11xxx
	ParamInvalid = BeError{"11001", errors.New("Invalid param")}
	TypeInvalid  = BeError{"11002", errors.New("Invalid type")}

	///////////////////////
	// Handler 12xxx
	///////////////////////

	///////////////////////
	// Service 14xxx
	///////////////////////
	ServiceCommon = BeError{"14000", errors.New("Service common error")}

	///////////////////////
	// Repository 16xxx
	///////////////////////
	RepositoryCommon = BeError{"16000", errors.New("Repository common error")}
	RecordNotfound   = BeError{"16004", errors.New("Record not found in repository")}
)
