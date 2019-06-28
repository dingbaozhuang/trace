package models

type HTTPResponse struct {
	Code int         `json:"code"`
	Err  string      `json:"err"`
	Data interface{} `json:"data"`
}

type RPCRequest struct {
	IPAddr string `json:"ip_addr"`
	Msg    string `json:"msg"`
}

type RPCResponse struct {
	Code int    `json:"code"`
	Err  string `json:"err"`
}
