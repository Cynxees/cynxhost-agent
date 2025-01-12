package response

import "cynxhostagent/internal/model/response/responsecode"

type APIResponse struct {
	Code     responsecode.ResponseCode `json:"code"`
	CodeName string                    `json:"codename"`
	Data     interface{}               `json:"data,omitempty"`  // Optional for success responses
	Error    string                    `json:"error,omitempty"` // Optional for error responses
}
