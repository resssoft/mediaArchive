package model

type RequestError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

type Response struct {
	Error    string      `json:"error,omitempty"`
	Total    int         `json:"total,omitempty"`
	Count    int         `json:"count,omitempty"`
	NextPage bool        `json:"nextPage,omitempty"`
	Data     interface{} `json:"data"`
}
