package goengine

import (
	"fmt"
	"net/http"
	"encoding/json"
)

type Response struct{
	w http.ResponseWriter
	template *Template
}

func NewResponse(w http.ResponseWriter) (*Response){
	response := &Response{}
	response.w = w
	return response
}

func (response *Response) SetHttpStatus(code int){
	response.w.WriteHeader(code)
}

func (response *Response) SetHeader(key string, value string){
	response.w.Header().Set(key, value)
}

func (response *Response) RenderText(text string){
	response.w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(response.w, text)
}

func (response *Response) RenderJson(data interface{}){
	bs, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	response.w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(response.w, string(bs))
}

func (response *Response) RenderHtml(tpl string, data interface{}){
	response.w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := response.template.AllTemplates.Lookup(tpl)
	if t != nil {
		t.Execute(response.w, data)
		return
	}
	fmt.Println("template not found", tpl)
}
