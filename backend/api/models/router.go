package models

import "go.mongodb.org/mongo-driver/bson/primitive"

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

type ImportParams struct {
	ServiceProcessing bool   `json:"serviceProcessing,omitempty"`
	DownloadToLocal   bool   `json:"downloadToLocal,omitempty"`
	FileName          string `json:"fileName,omitempty"`
}

type Filter map[string]interface{}

func (f Filter) ToPrimitive() []primitive.E {
	filterPrimitive := make([]primitive.E, 0)
	for k, v := range f {
		filterPrimitive = append(filterPrimitive, primitive.E{Key: k, Value: v})
	}
	return filterPrimitive
}
