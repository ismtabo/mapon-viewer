package dao

// MaponError represents a Mapon Error Response.
type MaponError struct {
	Err struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	} `json:"error"`
}
