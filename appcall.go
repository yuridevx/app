package app

import "github.com/yuridevx/app/options"

type AppCall struct {
	CallType options.CallType
}

func (c *AppCall) GetCallType() options.CallType {
	return c.CallType
}

func (c *AppCall) IsCallType(callType ...options.CallType) bool {
	for _, v := range callType {
		if v == c.CallType {
			return true
		}
	}
	return false
}

func (c *AppCall) GetHandler() uintptr {
	return 0
}

func (c *AppCall) GetComponentDefinition() interface{} {
	return nil
}

var _ options.Call = (*AppCall)(nil)
