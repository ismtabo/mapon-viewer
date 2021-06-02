package dao

type MaponError struct {
	Err struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	} `json:"error"`
}
