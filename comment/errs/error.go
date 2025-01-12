package errs

import "github.com/pkg/errors"

var ParamErr = errors.New("参数错误")

func NewParamErr(val string) error {
	return errors.Wrap(ParamErr, val)
}
