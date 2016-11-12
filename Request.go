package engine

import (
	"net/http"
	"strings"
)

type Request struct{
	IsMultiParamsParsed bool
	IsFormParamsParsed bool
	r *http.Request
}

func NewRequest(r *http.Request) (*Request){
	request := &Request{}
	request.r = r
	return request
}

func (request *Request)PostParamByName(name string) (value string){
	contentType := request.r.Header.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") == true {
		if request.IsMultiParamsParsed == false {
			request.r.ParseMultipartForm(32 << 20)
			request.IsMultiParamsParsed = true
		}
		
		if request.r.MultipartForm != nil {
			tmp := request.r.MultipartForm.Value[name]
			//WARNING only return the first value default
			if len(tmp) > 0 {
				return tmp[0]
			}
		}
	} else {
		if request.IsFormParamsParsed == false {
			request.r.ParseForm()
		}
		return request.r.PostFormValue(name)
	}
	return ""
}

func (request *Request)GetParamByName(name string) (value string){
	if request.IsFormParamsParsed == false {
		request.r.ParseForm()
	}
	
	return request.r.Form.Get(name)
}
