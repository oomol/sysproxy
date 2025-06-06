package sysproxy

import "errors"

var (
	ErrHttpsNotSupport = errors.New("os not support https proxy")
)
