package controllers

import (
	"net/http"
	"github.com/go-openapi/runtime"
)

// callback response
type CbResponse struct {
	cb	func(http.ResponseWriter, runtime.Producer)
}

func NewResponse(cb	func(http.ResponseWriter, runtime.Producer)) *CbResponse {
	return &CbResponse{cb:cb}
}

func (o *CbResponse) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	o.cb(rw, producer)
}

// simple response
type Response struct {
	StatusCode int
	Payload interface{}
}

func NewSuccess() *Response {
	resp := &Response{
		StatusCode: 200,
	}
	resp.Success()
	return resp
}

func (r *Response) Success() *Response {
	r.Payload = map[string]interface{}{
		"code": "success",
		"data": struct{}{},
	}
	return r
}

func (r *Response) Page(limit, offset, total int64, items interface{}) *Response {
	r.Payload = map[string]interface{}{
		"code": "success",
		"data": map[string]interface{}{
			"items": items,
			"limit": limit,
			"offset": offset,
			"total": total,
		},
	}
	return r
}

func (r *Response) List(limit, items, cursor interface{}) *Response {
	r.Payload = map[string]interface{}{
		"code": "success",
		"data": map[string]interface{}{
			"limit": limit,
			"cursor": cursor,
			"items": items,
		},
	}
	return r
}

func (r *Response) Item(name string, item interface{}) *Response {
	r.Payload = map[string]interface{}{
		"code": "success",
		"data": map[string]interface{}{
			name: item,
		},
	}
	return r
}

func (r *Response) SetData(data interface{}) *Response {
	r.Payload = map[string]interface{}{
		"code": "success",
		"data": data,
	}
	return r
}

func (r *Response) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(r.StatusCode)
	if r.Payload != nil {
		payload := r.Payload
		if err := producer.Produce(rw, payload); err != nil {
			// write tcp 10.42.14.3:3009->10.42.147.247:35846: write: broken pipe
			//log.Error(err.Error())
		}
	}
}
