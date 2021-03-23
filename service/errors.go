package service

import "errors"

//Default stream errors
var (
	Success             = "success"
	ErrorTaskNotFound   = errors.New("task not found")
	ErrorTaskExisit     = errors.New("task already exisit")
	ErrorTaskUpdateFail = errors.New("task update fail")

	ErrorClientNotFound = errors.New("client not found")
	ErrorClientExisit   = errors.New("client already exisit")
)
