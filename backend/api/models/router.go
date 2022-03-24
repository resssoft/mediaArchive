package models

import (
	"github.com/resssoft/mediaArchive/pkg/requestFilter"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

type DataFilter struct {
	Data []primitive.E
}

func (f *DataFilter) Append(field string, value interface{}) *DataFilter {
	f.Data = append(f.Data, primitive.E{Key: field, Value: value})
	return f
}

func (f *DataFilter) FromRequest(params map[string]string, data []byte) *DataFilter {
	reqfilters, _ := requestFilter.BuildFilter(params, data)
	f.AppendFromRequestFilter(reqfilters)
	return f
}

func (f *DataFilter) AppendFromRequestFilter(filter requestFilter.Filter) *DataFilter {
	if f.Data == nil {
		f.Data = []primitive.E{}
	}
	for _, reqFilter := range filter.Filters {
		var key string
		var val interface{}
		for k, v := range reqFilter.Data {
			key = k
			val = v
			break
		}
		log.Info().Interface("append filter by list", reqFilter).Send()
		f.Data = append(f.Data, primitive.E{Key: key, Value: val})
	}
	return f
}
